#include <BLEServerSetup.h>

NimBLEServer *bServer;
NimBLECharacteristic *wifiConnStatusCharacteristic;
NimBLECharacteristic *uploadProgressCharacteristic;
NimBLECharacteristic *uploadStatusCharacteristic;

bool isConnected = false;
//String wifiStatus = "notConnected";

void setupBLE() {
  Serial.println("Setting up bluetooth services");
  

  NimBLEDevice::init("Heartlink 01");
  bServer = NimBLEDevice::createServer();
  bServer->setCallbacks(new ServerConnectionCallbacks());

  // Configure Advertising services and characteristics
  NimBLEService *wifiService = bServer->createService(WIFI_SERVICE);
  NimBLEService *recordService = bServer->createService(RECORDING_SERVICE);
  NimBLEService *patientInfoService = bServer->createService(PATIENT_INFO_SERVICE);

  NimBLEService *uploadService = bServer->createService(UPLOAD_SERVICE);


  NimBLECharacteristic *connectWifiCharacteristic = wifiService->createCharacteristic(
    WIFI_CREDS_CHARACTERISITIC, NIMBLE_PROPERTY::WRITE
  );
  wifiConnStatusCharacteristic = wifiService->createCharacteristic(
    WIFI_CONNECTION_STATUS_CHARACTERISTIC, NIMBLE_PROPERTY::NOTIFY
  );
  NimBLECharacteristic *recordCharacteristic = recordService->createCharacteristic(
    RECORDING_CHARACTERISTIC, NIMBLE_PROPERTY::WRITE
  );
  NimBLECharacteristic *patientInfoCharacteristic = patientInfoService->createCharacteristic(
    PATIENT_INFO_CHARACTERISTIC, NIMBLE_PROPERTY::WRITE
  );
  uploadProgressCharacteristic = uploadService->createCharacteristic(
    UPLOAD_PROGRESS_CHARACTERISTIC, NIMBLE_PROPERTY::NOTIFY
  );
  uploadStatusCharacteristic = uploadService->createCharacteristic(
    UPLOAD_STATUS_CHARACTERISTIC, NIMBLE_PROPERTY::NOTIFY
  );

  connectWifiCharacteristic->setCallbacks(new WifiCallbacks());
  wifiConnStatusCharacteristic->setCallbacks(new WifiCallbacks());
  recordCharacteristic->setCallbacks(new RecordingCallbacks());
  patientInfoCharacteristic->setCallbacks(new PatientInfoCallbacks());

  uploadProgressCharacteristic->setCallbacks(new UploadCallbacks());
  uploadStatusCharacteristic->setCallbacks(new UploadCallbacks());

  wifiService->start();
  recordService->start();
  patientInfoService->start();

  uploadService->start();
  
  // Advertising device config
  NimBLEAdvertising *pAdvertising = NimBLEDevice::getAdvertising();
  pAdvertising->addServiceUUID(WIFI_SERVICE);
  pAdvertising->addServiceUUID(RECORDING_SERVICE);
  pAdvertising->addServiceUUID(UPLOAD_SERVICE);
  pAdvertising->start();
  
}