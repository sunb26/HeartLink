#include "i2sMicrophone.h"

File file;

const int headerSize = 44;

volatile bool isRecordingComplete = false;


void i2sInit() {
    i2s_config_t i2s_config = {
        .mode = (i2s_mode_t)(I2S_MODE_MASTER | I2S_MODE_RX),
        .sample_rate = I2S_SAMPLE_RATE,
        .bits_per_sample = i2s_bits_per_sample_t(I2S_SAMPLE_BITS),
        .channel_format = I2S_CHANNEL_FMT_ONLY_LEFT,
        .communication_format = I2S_COMM_FORMAT_I2S,  // Update this line
        .intr_alloc_flags = 0,
        .dma_buf_count = 4,
        .dma_buf_len = 512,
        .use_apll = 1
    };
 
    i2s_driver_install(I2S_PORT, &i2s_config, 0, NULL);
 
    const i2s_pin_config_t pin_config = {
        .bck_io_num = I2S_SCK,
        .ws_io_num = I2S_WS,
        .data_out_num = -1,
        .data_in_num = I2S_SD
    };
 
    i2s_set_pin(I2S_PORT, &pin_config);
}
 
 
void wavHeader(byte *header, int wavSize) {
  header[0] = 'R'; header[1] = 'I'; header[2] = 'F'; header[3] = 'F';
  unsigned int fileSize = wavSize + headerSize - 8;
  header[4] = (byte)(fileSize & 0xFF);
  header[5] = (byte)((fileSize >> 8) & 0xFF);
  header[6] = (byte)((fileSize >> 16) & 0xFF);
  header[7] = (byte)((fileSize >> 24) & 0xFF);
  header[8] = 'W'; header[9] = 'A'; header[10] = 'V'; header[11] = 'E';
  header[12] = 'f'; header[13] = 'm'; header[14] = 't'; header[15] = ' ';
  header[16] = 16; header[17] = 0; header[18] = 0; header[19] = 0;
  header[20] = 1; header[21] = 0;
  header[22] = I2S_CHANNEL_NUM; header[23] = 0;
  header[24] = (byte)(I2S_SAMPLE_RATE & 0xFF);
  header[25] = (byte)((I2S_SAMPLE_RATE >> 8) & 0xFF);
  header[26] = (byte)((I2S_SAMPLE_RATE >> 16) & 0xFF);
  header[27] = (byte)((I2S_SAMPLE_RATE >> 24) & 0xFF);
  unsigned int byteRate = I2S_SAMPLE_RATE * I2S_CHANNEL_NUM * I2S_SAMPLE_BITS / 8;
  header[28] = (byte)(byteRate & 0xFF);
  header[29] = (byte)((byteRate >> 8) & 0xFF);
  header[30] = (byte)((byteRate >> 16) & 0xFF);
  header[31] = (byte)((byteRate >> 24) & 0xFF);
  header[32] = I2S_CHANNEL_NUM * I2S_SAMPLE_BITS / 8;
  header[33] = 0;
  header[34] = I2S_SAMPLE_BITS;
  header[35] = 0;
  header[36] = 'd'; header[37] = 'a'; header[38] = 't'; header[39] = 'a';
  header[40] = (byte)(wavSize & 0xFF);
  header[41] = (byte)((wavSize >> 8) & 0xFF);
  header[42] = (byte)((wavSize >> 16) & 0xFF);
  header[43] = (byte)((wavSize >> 24) & 0xFF);
}
 
void i2s_adc(void *arg) {
  size_t bytes_read;
  char *i2s_read_buff = (char *)calloc(I2S_READ_LEN, sizeof(char));
 

  file = SD.open(filename, FILE_WRITE);
  if (!file) {
    Serial.println(F("Failed to open file for writing"));
    free(i2s_read_buff);
    vTaskDelete(NULL);
    return;
  }
 
  byte header[headerSize];
  wavHeader(header, FLASH_RECORD_SIZE);
  file.write(header, headerSize);
 
  int totalBytesWritten = 0;
    
  Serial.println(F("*** Recording Start ***"));
  digitalWrite(RED_LED_PIN, HIGH);
  
  while (totalBytesWritten < FLASH_RECORD_SIZE && startStop == "start") { //might be able to add an && in this if statement for stop/start
    i2s_read(I2S_PORT, (void *)i2s_read_buff, I2S_READ_LEN, &bytes_read, portMAX_DELAY);
    file.write((const uint8_t *)i2s_read_buff, bytes_read);
    totalBytesWritten += bytes_read;
    Serial.printf("Progress: %d%%\n", (totalBytesWritten * 100) / FLASH_RECORD_SIZE);
  }

  fileSize = totalBytesWritten;
  file.flush();
  file.close();
  free(i2s_read_buff);
  digitalWrite(RED_LED_PIN, LOW);
  startStop = "stop";
  Serial.println(F("*** Recording Complete ***"));
  isRecordingComplete = true; // Recording is now complete
  
  vTaskDelete(NULL);
}
