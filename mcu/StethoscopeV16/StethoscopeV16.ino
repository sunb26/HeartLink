#include "Config.h"
#include "InitializeSD.h"
#include "ServerUpload.h"
#include "i2sMicrophone.h"
#include "BLEServerSetup.h"
#include "BLEServiceCallbacks.h"

bool firstRun = true;
int fileSize = 0;

void setup(void) {
  
  if (firstRun == true){
    Serial.begin(115200);
  
    pinMode(GREEN_LED_PIN, OUTPUT);
    pinMode(WHITE_LED_PIN, OUTPUT);
    pinMode(RED_LED_PIN, OUTPUT);
    pinMode(BLUE_LED_PIN, OUTPUT);

    analogWrite(GREEN_LED_PIN, 50);
    analogWrite(WHITE_LED_PIN, 0);
    analogWrite(BLUE_LED_PIN, 0);

    //digitalWrite(GREEN_LED_PIN, HIGH);
    //digitalWrite(WHITE_LED_PIN, LOW); 
    digitalWrite(RED_LED_PIN, LOW);
    //digitalWrite(BLUE_LED_PIN, LOW);
    
    setupBLE();
    initializeSD(); //try to put this first before the if statement

    while (ssid == NULL){
      delay(10);
    }
    initializeWifi();

    firstRun = false;

    Serial.println(startStop);
    while(startStop == "stop"){
      delay (10);
      }
    // Start I2S ADC task
    i2sInit();
  }
  
  /*else{
  // Start I2S ADC task
  //  i2sInit();
  }*/

  xTaskCreate(i2s_adc, "i2s_adc", 1024 * 4, NULL, 1, NULL);
}
 
void loop(void) {

  if (isRecordingComplete == true){ 
    if (fileUploaded == false){
      uploadFileToServer(filename);
      fileUploaded = true;
    }

    if (startStop == "start"){
      isRecordingComplete = false;
      fileUploaded = false;
      setup();
    }
  }
  wifiStatus();
    
  delay(20);  //allow the cpu to switch to other tasks
}
