# CHANGELOG

## 1.18.0 (2025-03-22)
- add handler logging 
- add listPatients endpoint

## 1.17.1 (2025-03-22)
- update database schema patient table with last_updated column

## 1.17.0 (2025-03-13)
- add server endpoint to upload .wav file with access token to Firebase 

## 1.16.2 (2025-03-10)
- fix bug in iOS app where wifi doesn't disconnect if bluetooth disconnects

## 1.16.1 (2025-03-06)
- update patient table text colours and header titles

## 1.16.0 (2025-03-05)
- add comments submission box on website
- add indication for patients with unviewed recordings to user web homepage

## 1.15.1 (2025-03-01)
- update iOS app color scheme

## 1.15.0 (2025-02-28)
- DB schema v2
- add audio player to website

## 1.14.0 (2025-02-21)
- server now hosted on google cloud via cloud run 

## 1.13.0 (2025-02-21)
- website deployed on Vercel

## 1.12.0 (2025-02-21)
- add patient list page
  - protect URL from un-authenticated users
  - data is currently mocked (hardcoded)
  - column sorting
  - global filtering
  - pagination
- add patient page
  - recording table
  - placeholder patient info
  - placeholder patient profile avatar

## 1.11.0 (2025-02-19)
- add initial server code in Golang (runs locally only)
- add website landing page
  - header
  - hero image
  - about section
  - getting started section
- add Clerk auth sign in option

## 1.10.1 (2025-01-28)
- add patient id sending to MCU across BLE (sent on each recording)
- add BLE service to app to subscribe to wifi connection status on MCU

## 1.9.1 (2025-01-27)
- fix logo colour

## 1.9.0 (2025-01-27)
- add confirmation dialog to app before submitting or deleting a recording

## 1.8.0 (2025-01-26)
- add submit and delete buttons for recordings
- separate recording history list from User model
- call recording history list endpoint upon appearance of homepage
- update logo

## 1.7.1 (2025-01-22)
- fix bug where audio continues to play if not paused before returning to homepage

## 1.7.0 (2025-01-22)
- add audio track seeking

## 1.6.1 (2025-01-22)
- fix bug where audio was not played if phone ringer was off

## 1.6.0 (2025-01-22)
- install Firebase SDK
- add audio streaming from Firebase Cloud Storage

## 1.5.0 (2025-01-21)
- add individual recording data view
- TODO (Ben Sun): add audio streaming to this view

## 1.4.0 (2025-01-20)
- remove hardcoded recording history
- replaced with pulling data from endpoint

## 1.3.0 (2025-01-18)
- set up wifi credentials login uses phone app now
- start/stop function uses phone app now
- app login page pulls from an endpoint (still no validation yet)
- recording page auto stops after 15 seconds
- TODO (bentomka): MCU issue with repeat recordings

## v1.2.0 (2025-01-14)
- setup atlasgo for SQL db schema management

## v1.1.0 (2025-01-14)
- add function to send Wifi Credentials to MCU via BLE
- add function to trigger start/stop recording via BLE to MCU

## v1.0.0 (2024-12-02)
- remove BLE audio streaming from MCU
- add file upload to server (localhost) to MCU
- add SD card integration to MCU

## v0.5.0 (2024-11-19)
- added audio sign .wav constructor feature

## v0.4.1 (2024-11-18)
- updated recording page to add countdown feature 

## v0.4.0 (2024-11-10)
- added bluetooth connection and view

## v0.3.0 (2024-10-31)
- added homepage
- added recording page
- refactored login page
- add MainNavView to handle navigation between existing pages
  
## v0.2.0 (2024-10-14)
- added login view (no functionality yet)

## v0.1.0 (2024-10-14)
- Initial iOS setup
