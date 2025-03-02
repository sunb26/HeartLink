//
//  SendWifiCreds.swift
//  HeartLink
//
//  Created by Ben Sun on 2025-01-11.
//

import SwiftUI

struct SendWifiCreds: View {
    @Environment(\.dismiss) var dismiss // For closing the pop-up
    @State private var network: String = ""
    @State private var password: String = ""
    @ObservedObject var bluetoothManager: BluetoothManager

    var body: some View {
        VStack {
            HStack {
                Text("Connect to WiFi")
                    .font(.title)
                    .bold()
                Image(systemName: "wifi")
                    .resizable()
                    .frame(width: 22, height: 22)
            }
            .padding(.bottom, 40)

            TextField("Network Name", text: $network)
                .padding(20)
                .background(Color.white)
                .cornerRadius(10)
                .shadow(color: .gray, radius: 5, x: 0, y: 2)

            SecureField("Password", text: $password)
                .padding(20)
                .background(Color.white)
                .cornerRadius(10)
                .shadow(color: .gray, radius: 5, x: 0, y: 2)
                .padding(.top, 12)

            if bluetoothManager.wifiConnStatus == "connected" {
                Button(action: {
                    dismiss()
                }) {
                    Text("Close")
                        .padding()
                        .frame(maxWidth: .infinity)
                        .background(Color("custom-red"))
                        .foregroundColor(.white)
                        .cornerRadius(8)
                        .shadow(color: .gray, radius: 5, x: 0, y: 2)
                }
                .padding(.top, 90)
            } else {
                Button(action: {
                    sendWifiCreds(network: network, password: password)
                }) {
                    if bluetoothManager.wifiConnStatus == "notConnected" {
                        Text("Connect")
                            .padding()
                            .frame(maxWidth: .infinity)
                            .background(Color("custom-red"))
                            .foregroundColor(.white)
                            .cornerRadius(8)
                            .shadow(color: .gray, radius: 5, x: 0, y: 2)
                    } else if bluetoothManager.wifiConnStatus == "connecting" {
                        ProgressView()
                    }
                }
                .padding(.top, 90)
            }
        }.padding(20)
    }

    func sendWifiCreds(network: String, password: String) {
        let nameLength = String(network.count)
        let data = (nameLength + "&" + network + password).data(using: .utf8)!
        guard let char = bluetoothManager.wifiCredsCharacteristic else {
            print("Could not find wifiCreds characteristic")
            return
        }
        bluetoothManager.mcuPeripheral?.writeValue(data, for: char, type: .withResponse)
    }
}

#Preview {
    @Previewable @StateObject var btmanager = BluetoothManager()

    SendWifiCreds(bluetoothManager: btmanager)
}
