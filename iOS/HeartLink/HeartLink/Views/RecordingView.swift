//
//  RecordingView.swift
//  HeartLink
//
//  Created by Ben Sun on 2024-10-29.
//  Edited by Matt Wilker on 2024-11-16

import SwiftUI

struct RecordingView: View {
    @State var startRecording: Bool = false
    @State var isRecording: Bool = false
    @State var countdown: Int8 = 0
    @State var recordingDuration: Int8 = 15
    @State var timer = Timer.publish(every: 1.0, on: .main, in: .common).autoconnect()
    @ObservedObject var bluetoothManager: BluetoothManager
    @Binding var patient: User

    func toggleRecording() {
        let data = (startRecording ? "start" : "stop").data(using: .utf8)!
        guard let char = bluetoothManager.recordingCharacteristic else {
            print("Could not find recording characteristic")
            return
        }
        bluetoothManager.mcuPeripheral?.writeValue(data, for: char, type: .withResponse)
    }

    var body: some View {
        if bluetoothManager.isConnected && bluetoothManager.wifiConnStatus == "connected" {
            ZStack {
                VStack {
                    Text(startRecording ? "Recording..." : "Ready to Record").font(.system(size: 42, weight: .bold))
                        .frame(maxWidth: .infinity, alignment: .center)
                        .frame(height: 100)

                    Text(countdown <= 0 ? " " : "Starts in: \(countdown)")
                        .font(.system(size: 34, weight: .bold))
                        .frame(maxWidth: .infinity, alignment: .center)
                        .frame(height: 400, alignment: .top)

                    Button(action: {
                        recordingDuration = 15
                        if startRecording {
                            startRecording = false
                            countdown = -1
                            toggleRecording()
                        } else {
                            startRecording = true
                            countdown = 3
                        }
                    }) {
                        if !startRecording { // ready to record page
                            Image(systemName: "record.circle.fill")
                                .resizable()
                                .scaledToFit()
                                .frame(width: 80, height: 80)
                                .foregroundColor(.red)
                        } else { // recording page
                            Image(systemName: "stop.circle")
                                .resizable()
                                .scaledToFit()
                                .frame(width: 80, height: 80)
                                .foregroundColor(.red)
                        }
                    }
                }
                .onReceive(timer) { _ in // decrement timer every second
                    guard startRecording else { return }

                    if countdown > 0 {
                        countdown -= 1
                        if countdown == 0 {
                            guard let char = bluetoothManager.patientInfoCharacteristic else {
                                print("Could not find patientInfo characteristic")
                                return
                            }
                            let patId = "\(patient.patientId)".data(using: .utf8)!
                            bluetoothManager.mcuPeripheral?.writeValue(patId, for: char, type: .withResponse)
                            
                            toggleRecording()
                        }
                    } else {
                        recordingDuration -= 1
                        print(recordingDuration)
                        if recordingDuration <= 0 {
                            print("stopping recording...")
                            startRecording = false
                            toggleRecording()
                            recordingDuration = 15
                        }
                    }
                }
            }
        } else {
            if !bluetoothManager.isConnected && bluetoothManager.wifiConnStatus != "connected" {
                Text("Please connect to the Bluetooth Device and Wifi before recording (found in the main menu)")
                    .font(.title3)
                    .padding(10)
                    .multilineTextAlignment(.center)
            } else if !bluetoothManager.isConnected {
                Text("Please connect to the Bluetooth Device before recording (found in the main menu)")
                    .font(.title3)
                    .padding(10)
                    .multilineTextAlignment(.center)
            } else {
                Text("Please connect device to Wifi before recording (found in the main menu)")
                    .font(.title3)
                    .padding(10)
                    .multilineTextAlignment(.center)
            }
        }
    }
}

#Preview {
    @Previewable var bt = BluetoothManager()
    @Previewable @State var patient: User = User(email: "test", patientId: 1, physicianId: 1)
    RecordingView(bluetoothManager: bt, patient: $patient)
}
