#include "BLEServiceCallbacks.h"

const char *ssid = NULL;
const char *password = NULL;
const char *patientID = NULL;
String startStop = "stop";

//char filename[64] = "/heartlink_0.wav";
char filename[64];

void ServerConnectionCallbacks::onConnect(NimBLEServer* pServer, NimBLEConnInfo& connInfo) {
    Serial.println("BLUETOOTH CONNECTED");
    //digitalWrite(BLUE_LED_PIN, HIGH);
    analogWrite(BLUE_LED_PIN, 100);

}

void ServerConnectionCallbacks::onDisconnect(NimBLEServer* pServer, NimBLEConnInfo& connInfo, int reason) {
    Serial.println("BLUETOOTH DISCONNECTED");
    NimBLEDevice::startAdvertising();
    analogWrite(BLUE_LED_PIN, 0);
    initializeWifi();
}

void WifiCallbacks::onWrite(NimBLECharacteristic *pCharacteristic, NimBLEConnInfo& connInfo) {
    String value = pCharacteristic->getValue();
    // Serial.println(value);
    // TODO: implement code to connect to WiFi
    
    int delimiterIndex = value.indexOf('&');

    int ssidLength = value.substring(0, delimiterIndex).toInt();
    String combined = value.substring(delimiterIndex + 1);

    static char parsedSSID[32];  // Adjust size to the max expected SSID length
    static char parsedPassword[64];  // Adjust size to the max expected password length

    combined.substring(0, ssidLength).toCharArray(parsedSSID, sizeof(parsedSSID));
    combined.substring(ssidLength).toCharArray(parsedPassword, sizeof(parsedPassword));

    ssid = parsedSSID;
    password = parsedPassword;

}

void RecordingCallbacks::onWrite(NimBLECharacteristic *pCharacteristic, NimBLEConnInfo& connInfo) {
    startStop = pCharacteristic->getValue();
    Serial.println(startStop);
}

void PatientInfoCallbacks::onWrite(NimBLECharacteristic *pCharacteristic, NimBLEConnInfo& connInfo) {
    String value = pCharacteristic->getValue();
    Serial.println("Patient Id:");
    Serial.println(value);  

    
    snprintf(filename, sizeof(filename), "/heartlink_%s.wav", value.c_str());
    Serial.print("Filename: ");
    Serial.println(filename);
}

void UploadCallbacks::onWrite(NimBLECharacteristic *pCharacteristic, NimBLEConnInfo& connInfo) {
}
