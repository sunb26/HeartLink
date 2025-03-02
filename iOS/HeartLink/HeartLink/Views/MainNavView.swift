//
//  MainNavView.swift
//  HeartLink
//
//  Created by Ben Sun on 2024-10-29.
//

import SwiftUI

// Page actions for navigation
enum PageActions: Hashable {
    case login
    case home
    case recording
}

struct MainNavView: View {
    @State var path: [PageActions] = [.login]
    @StateObject var btmanager = BluetoothManager()
    @State var patient: User = User(email: "", patientId: 0, physicianId: 0)
    @State var recordingData: RecordingData = RecordingData(id: 0, date: "0000-00-00", viewStatus: "notSubmitted", comments: "", fileURL: "")

    var body: some View {
        NavigationStack(path: $path) {
            TabView {
                Tab("Home", systemImage: "house.fill") {
                    HomeView(path: $path, patient: $patient, recordingData: $recordingData, bluetoothManager: btmanager)
                }
                Tab("Record", systemImage: "record.circle.fill") {
                    RecordingView(bluetoothManager: btmanager, patient: $patient)
                }
            }
            .tabViewStyle(.sidebarAdaptable)
            .tint(Color("custom-red"))
            .edgesIgnoringSafeArea(.all)
            .navigationDestination(for: PageActions.self) { action in
                switch action {
                case .login:
                    LoginView(path: $path, patient: $patient)
                        .navigationBarBackButtonHidden(true)
                case .home:
                    HomeView(path: $path, patient: $patient, recordingData: $recordingData, bluetoothManager: btmanager)
                case .recording:
                    RecLogView(path: $path, recording: $recordingData)
                }
            }
        }
    }
}

#Preview {
    MainNavView()
}
