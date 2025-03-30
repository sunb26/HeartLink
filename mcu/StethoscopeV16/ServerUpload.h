#include "Config.h"
#include "InitializeSD.h"

#include <WiFi.h>

#include <WiFiClientSecure.h>

#include "BLEServiceCallbacks.h"
#include <BLEServerSetup.h>


extern volatile bool fileUploaded;

void initializeWifi();
void uploadFileToServer(const char* path);
void wifiStatus();
