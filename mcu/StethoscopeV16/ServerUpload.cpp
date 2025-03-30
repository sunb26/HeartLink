#include "ServerUpload.h"
 
volatile bool fileUploaded = false;
bool connected = false;


int statusFromResponse(WiFiClient &client) {
  size_t p = 0;
  size_t available = 0;
  // wait for first byte
  while (client.connected()) {
    if (client.available()) {
      available = 1;
      break;
    }
    delay(1);
    p++;
  }
  if (!available) {
 
    //Serial.println("not connected");
    //Serial.println(client.connected());
    return 0;
  }
  //Serial.printf("waited %d for %d available --->\n", p, available);
 
  int code = -1;
  size_t t = 0;
  char buf[200];
 
  while (client.connected() && available) {
    p = client.read((uint8_t *)buf, std::min(available, sizeof(buf) - 1));
    if (t == 0) {
      sscanf(buf, "HTTP/1.1 %*s %d", &code);
      //Serial.println("Code Result:");
      //Serial.println(code);
    }
    t += p;
    buf[p] = 0;
    Serial.print(buf);
    available = client.available();
  }
  //Serial.printf("<---\n%d total\n", t);
  return code;
}
 
 
void uploadFileToServer(const char* path) {
  // Check WiFi status first
  if (WiFi.status() != WL_CONNECTED) {
    //Serial.println(F("WiFi not connected. Attempting to reconnect..."));
    WiFi.reconnect();
   
    // Wait up to 5 seconds for reconnection
    int attempts = 0;
    while (WiFi.status() != WL_CONNECTED && attempts < 10) {
      delay(500);
      attempts++;
    }
   
    if (WiFi.status() != WL_CONNECTED) {
      //Serial.println(F("Failed to reconnect WiFi. Aborting upload."));
      return;
    }
  }
 
  if (!SD.exists(path)) {
    //Serial.println(F("File does not exist on SD card."));
    return;
  }
 
  // Open the file from the SD card
  File audioFile = SD.open(path);
  if (!audioFile) {
    //Serial.println(F("Failed to open file"));
    return;
  }
  
  WiFiClientSecure client;

  //Serial.println("URL: " + String(serverUrl));
 
  client.setInsecure();
 
  int connectionAttempts = 0;

  while (!client.connected() && connectionAttempts < 3) {
    connectionAttempts++;
    //Serial.printf("Connection attempt %d of %d\n", connectionAttempts, 3);
   
    client.stop();
    connected = client.connect(serverUrl, port,8000);
    if (!client.connected()) {
      //Serial.println("Connection failed, retrying...");

      client.flush();
      client.stop();
      delay(1000);  // Wait before retry
    }
  }
 
  if (client.connected()) {
    //Serial.println("Connected to server, uploading file...");
   
    // Your existing upload code
    client.println("POST https://heartlink-652851748566.northamerica-northeast2.run.app/UploadFilterRecording HTTP/1.1");
    client.print("Host: ");
    client.println(serverUrl);
 
    String boundary = "..--..--MyOwnBoundary" + String(random(0xEFFFFF) + 0x100000, HEX);
    client.println("Content-Type: multipart/form-data; boundary=" + boundary);
    client.println("Transfer-Encoding: chunked");
    client.println();
   
    // Rest of your upload code...
    auto filePart = [&](const String &name, File &file, const String &filename, const String &type) {
      String part = "\r\n--" + boundary + "\r\n";
      part += "Content-Disposition: form-data; name=\"";
      part += name;
      part += "\"; filename=\"";
      part += filename;
      part += "\"\r\nContent-Type: ";
      part += type;
      part += "\r\n\r\n";
 
      client.println(part.length(), HEX);
      client.println(part);
 
      const size_t BUFFER_SIZE = 4096;  // 8 KB buffer

      uint8_t* buf = (uint8_t*)malloc(BUFFER_SIZE);
      if (!buf) {
          Serial.println("Memory allocation failed!");
          return;
      }
      float increment = 4096.0/float(fileSize);
      //Serial.printf("Increment: %.2f%\n", increment);
      //Serial.printf("Total Bytes: %.2f%\n", float(totalBytesWritten));

      float progress = 0;
      size_t r = file.read(buf, BUFFER_SIZE);
      while (r > 0) {
          client.println(r, HEX);
          client.write(buf, r);
          client.println();
          //Serial.printf("wrote %d of %s, now position %d\n", r, file.name(), file.position());
          progress += increment;
          //Serial.printf("Upload progress: %.2f\n", progress);
          uploadProgressCharacteristic->setValue(String(progress)); //is this the right value for Benson
          uploadProgressCharacteristic->notify();
          r = file.read(buf, BUFFER_SIZE);
      }
      client.flush();
      free(buf);

    };
 
    auto endPart = [&]() {
      String part = "\r\n--" + boundary + "--\r\n";
      client.println(part.length(), HEX);
      client.println(part);
      client.println(0);
      client.println();
    };
 
    filePart("audioKey", audioFile, path + 1, "audio/wav");
    endPart();

    uploadProgressCharacteristic->setValue(String(1)); //is this the right value for Benson
    uploadProgressCharacteristic->notify();
 
    int status = statusFromResponse(client);
    Serial.printf("Upload status: %d\n", status);
   
    // Handle the response status
    if (status == 0) {
      uploadStatusCharacteristic->setValue("failed"); //is this the right value for Benson
      uploadStatusCharacteristic->notify();
      WiFi.reconnect();
    } else {
      uploadStatusCharacteristic->setValue("success"); //is this the right value for Benson
      uploadStatusCharacteristic->notify();
    }
 
    audioFile.close();
  }

  else {
    //Serial.println("Failed to connect to server after multiple attempts");
    audioFile.close();
  }
  client.flush();
  client.stop();
 
  // Set a flag to indicate completion
  fileUploaded = true;
}
 
void wifiStatus(){
  if (WiFi.status() == WL_CONNECTED){
    wifiConnStatusCharacteristic->setValue("connected");
    wifiConnStatusCharacteristic->notify();
  }
 
  else{
    wifiConnStatusCharacteristic->setValue("connecting");
    wifiConnStatusCharacteristic->notify();
    initializeWifi();
  }
 
}
 
void initializeWifi(){
    while (ssid == NULL){
      delay(10);
    }
    WiFi.mode(WIFI_STA);
    WiFi.begin(ssid, password);
    Serial.print(F("Connecting..."));
    wifiConnStatusCharacteristic->setValue("connecting"); //is this the right value for Benson
    wifiConnStatusCharacteristic->notify();
 
    // Wait for connection
    uint8_t i = 0;
    while (WiFi.status() != WL_CONNECTED && i++ < 10) {  //wait 5 seconds
        delay(500);
        analogWrite(WHITE_LED_PIN, 0);
    }
    if (i == 11) { // maybe add some feedback to user app about this
        //Serial.print(F("Could not connect to"));
        //DBG_OUTPUT_PORT.println(ssid);
        analogWrite(WHITE_LED_PIN, 0);
        wifiConnStatusCharacteristic->setValue("notConnected"); //is this the right value for Benson
        wifiConnStatusCharacteristic->notify();
        ssid = NULL;
        password = NULL;
        initializeWifi();
    }
    //Serial.print(F("Connected!"));
    analogWrite(WHITE_LED_PIN, 100);
 
    wifiStatus();
}
 
 