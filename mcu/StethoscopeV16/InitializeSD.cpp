#include "InitializeSD.h"

bool hasSD = false;

SPIClass spi = SPIClass(SPI);


void initializeSD(){
  
    spi.begin(SCK, MISO, MOSI, CS);
    if (SD.begin(CS, spi, 80000000)) {
        Serial.println(F("SD Card initialized."));
        hasSD = true;
    }
    else{
      Serial.println(F("Failed to initialize SD Card."));
    }
}

