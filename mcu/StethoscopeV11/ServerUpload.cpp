#include "ServerUpload.h"
#include "BLEServiceCallbacks.h"

File uploadFile;


const char *host = "esp32sd";

WebServer server(80);

void returnOK() {
  server.send(200, "text/plain", "");
}

void returnFail(String msg) {
  server.send(500, "text/plain", msg + "\r\n");
}

bool loadFromSdCard(String path) {
  String dataType = "text/plain";
  if (path.endsWith("/")) {
    path += "index.htm";
  }
 
 
  File dataFile = SD.open(path.c_str());
  if (dataFile.isDirectory()) {
    path += "/index.htm";
    dataType = "text/html";
    dataFile = SD.open(path.c_str());
  }
 
  if (!dataFile) {
    return false;
  }
 
  if (server.hasArg("download")) {
    dataType = "application/octet-stream";
  }
 
  if (server.streamFile(dataFile, dataType) != dataFile.size()) {
    DBG_OUTPUT_PORT.println("Sent less data than expected!");
  } 
  
  dataFile.close();
  return true;
}

void handleFileUpload() {
  if (server.uri() != "/edit") {
    return;
  }
  HTTPUpload &upload = server.upload();
  if (upload.status == UPLOAD_FILE_START) {
    if (SD.exists((char *)upload.filename.c_str())) {
      SD.remove((char *)upload.filename.c_str());
    }
    uploadFile = SD.open(upload.filename.c_str(), FILE_WRITE);
    DBG_OUTPUT_PORT.print("Upload: START, filename: ");
    DBG_OUTPUT_PORT.println(upload.filename);
  } else if (upload.status == UPLOAD_FILE_WRITE) {
    if (uploadFile) {
      uploadFile.write(upload.buf, upload.currentSize);
    }
    DBG_OUTPUT_PORT.print("Upload: WRITE, Bytes: ");
    DBG_OUTPUT_PORT.println(upload.currentSize);
    
  } else if (upload.status == UPLOAD_FILE_END) {
    if (uploadFile) {
      uploadFile.close();
    }
    DBG_OUTPUT_PORT.print("Upload: END, Size: ");
    DBG_OUTPUT_PORT.println(upload.totalSize);
  }
}

void deleteRecursive(String path) {
  File file = SD.open((char *)path.c_str());
  if (!file.isDirectory()) {
    file.close();
    SD.remove((char *)path.c_str());
    return;
  }
 
  file.rewindDirectory();
  while (true) {
    File entry = file.openNextFile();
    if (!entry) {
      break;
    }
    String entryPath = path + "/" + entry.name();
    if (entry.isDirectory()) {
      entry.close();
      deleteRecursive(entryPath);
    } else {
      entry.close();
      SD.remove((char *)entryPath.c_str());
    }
    yield();
  }
 
  SD.rmdir((char *)path.c_str());
  file.close();
}

void handleDelete() {
  if (server.args() == 0) {
    return returnFail("BAD ARGS");
  }
  String path = server.arg(0);
  if (path == "/" || !SD.exists((char *)path.c_str())) {
    returnFail("BAD PATH");
    return;
  }
  deleteRecursive(path);
  returnOK();
}
 
void handleCreate() {
  if (server.args() == 0) {
    return returnFail("BAD ARGS");
  }
  String path = server.arg(0);
  if (path == "/" || SD.exists((char *)path.c_str())) {
    returnFail("BAD PATH");
    return;
  }
 
  if (path.indexOf('.') > 0) {
    File file = SD.open((char *)path.c_str(), FILE_WRITE);
    if (file) {
      file.write(0);
      file.close();
    }
  } else {
    SD.mkdir((char *)path.c_str());
  }
  returnOK();
}
 
void printDirectory() {
  if (!server.hasArg("dir")) {
    return returnFail("BAD ARGS");
  }
  String path = server.arg("dir");
  if (path != "/" && !SD.exists((char *)path.c_str())) {
    return returnFail("BAD PATH");
  }
  File dir = SD.open((char *)path.c_str());
  path = String();
  if (!dir.isDirectory()) {
    dir.close();
    return returnFail("NOT DIR");
  }
  dir.rewindDirectory();
  server.setContentLength(CONTENT_LENGTH_UNKNOWN);
  server.send(200, "text/json", "");
 
  server.sendContent("[");
  for (int cnt = 0; true; ++cnt) {
    File entry = dir.openNextFile();
    if (!entry) {
      break;
    }
 
    String output;
    if (cnt > 0) {
      output = ',';
    }
 
    output += "{\"type\":\"";
    output += (entry.isDirectory()) ? "dir" : "file";
    output += "\",\"name\":\"";
    output += entry.path();
    output += "\"";
    output += "}";
    server.sendContent(output);
    entry.close();
  }
  server.sendContent("]");
  dir.close();
}
 
void handleNotFound() {
  if (hasSD && loadFromSdCard(server.uri())) {
    return;
  }
  String message = "SDCARD Not Detected\n\n";
  message += "URI: ";
  message += server.uri();
  message += "\nMethod: ";
  message += (server.method() == HTTP_GET) ? "GET" : "POST";
  message += "\nArguments: ";
  message += server.args();
  message += "\n";
  for (uint8_t i = 0; i < server.args(); i++) {
    message += " NAME:" + server.argName(i) + "\n VALUE:" + server.arg(i) + "\n";
  }
  server.send(404, "text/plain", message);
  DBG_OUTPUT_PORT.print(message);
}


void initializeServer(){
    WiFi.mode(WIFI_STA);
    WiFi.begin(ssid, password);
    DBG_OUTPUT_PORT.print("Connecting to ");
    DBG_OUTPUT_PORT.println(ssid);
    
    // Wait for connection
    uint8_t i = 0;
    while (WiFi.status() != WL_CONNECTED && i++ < 20) {  //wait 10 seconds
        delay(500);
        digitalWrite(BLUE_LED_PIN, LOW);
    }
    if (i == 21) {
        DBG_OUTPUT_PORT.print("Could not connect to");
        DBG_OUTPUT_PORT.println(ssid);
        digitalWrite(BLUE_LED_PIN, LOW);
        while (1) {
        delay(500);
        }
    }
    DBG_OUTPUT_PORT.print("Connected! IP address: ");
    digitalWrite(BLUE_LED_PIN, HIGH);
    DBG_OUTPUT_PORT.println(WiFi.localIP());
    
    if (MDNS.begin(host)) {
        MDNS.addService("http", "tcp", 80);
        DBG_OUTPUT_PORT.println("MDNS responder started");
        DBG_OUTPUT_PORT.print("You can now connect to http://");
        DBG_OUTPUT_PORT.print(host);
        DBG_OUTPUT_PORT.println(".local");
    }
    
    server.on("/list", HTTP_GET, printDirectory);
    server.on("/edit", HTTP_DELETE, handleDelete);
    /////////////////////////////////////////////////Look into this part, maybe make a different function to post
    server.on("/edit", HTTP_PUT, handleCreate);
    server.on(
        "/edit", HTTP_POST,
        []() {
        returnOK();
        },
        handleFileUpload
    );
    server.onNotFound(handleNotFound);
    
    server.begin();
    DBG_OUTPUT_PORT.println("HTTP server started"); 
}
