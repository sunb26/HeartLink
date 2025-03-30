#ifndef BLE_SERVER_SETUP_H
#define BLE_SERVER_SETUP_H

#include <NimBLEDevice.h>
#include <BLEServiceCallbacks.h>

// Services
#define WIFI_SERVICE "5c96e1a0-4022-4310-816f-bcb7245bc802"
#define RECORDING_SERVICE "60ec2f71-22f2-4fc4-84f0-f8d3269e10c0"
#define PATIENT_INFO_SERVICE "a718bad4-f9b0-40c8-bd02-de0b1335aabb"

#define UPLOAD_SERVICE "fa680e3d-557d-4848-b6c9-8a9e3f149184"

// Characteristics
#define WIFI_CREDS_CHARACTERISITIC "a48ce354-6a1b-429d-aca5-1077627d5a25"
#define WIFI_CONNECTION_STATUS_CHARACTERISTIC "028807ff-751e-4798-a168-cb391c05288f"
#define RECORDING_CHARACTERISTIC "d5435c8c-392f-4e89-87be-89f9964db0e0"
#define PATIENT_INFO_CHARACTERISTIC "b69a6f81-c0fa-4aab-8bdf-84796a3f0aab"

#define UPLOAD_PROGRESS_CHARACTERISTIC "c39a4162-3362-485b-a799-133e32b3ac32"
#define UPLOAD_STATUS_CHARACTERISTIC "9ecc5e84-19c8-4120-8a0b-df0205b249ee"



void setupBLE();

extern NimBLECharacteristic *wifiConnStatusCharacteristic;

extern NimBLECharacteristic *uploadProgressCharacteristic;

extern NimBLECharacteristic *uploadStatusCharacteristic;

#endif
