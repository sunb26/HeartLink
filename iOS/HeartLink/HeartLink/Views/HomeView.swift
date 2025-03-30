//
//  HomeView.swift
//  HeartLink
//
//  Created by Ben Sun on 2024-10-29.
//

import SwiftUI

struct HomeView: View {
    @Binding var path: [PageActions]
    @Binding var isLoggedIn: Bool
    @Binding var patient: User
    @Binding var recordingData: RecordingData
    @State var btPopUp: Bool = false
    @State var recordingPopUp: Bool = false
    @State private var recordingsList: RecordingsList = RecordingsList(widgets: [])
    @ObservedObject var bluetoothManager: BluetoothManager

    var body: some View {
        ZStack {
            Color.white.ignoresSafeArea(.all)
            VStack {
                HStack {
                    Button(action: {
                        // Navigate to profile sidebar
                    }) {
                        Image(systemName: "person.circle")
                            .resizable()
                            .scaledToFit()
                            .frame(width: 30, height: 30)
                            .foregroundColor(.black)
                    }
                    .frame(maxWidth: 300, alignment: .leading)

                    Button(action: {
                        self.btPopUp = true
                    }) {
                        Image("bluetooth")
                            .resizable()
                            .scaledToFit()
                            .frame(width: 30, height: 30)
                            .foregroundColor(.black)
                    }
                }
                .frame(height: 80)

                Text("Home")
                    .font(.system(size: 34, weight: .bold))
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding(.horizontal)
                    .padding(.top, -4)

                VStack(alignment: .leading) {
                    Text("Recordings")
                        .font(.title)
                        .fontWeight(.bold)
                        .padding(.horizontal)

                    List(recordingsList.widgets) { recording in
                        HStack(spacing: 16) {
                            Image(systemName: "waveform")
                                .resizable()
                                .scaledToFit()
                                .frame(width: 40, height: 40)
                                .foregroundColor(Color("custom-red"))

                            VStack(alignment: .leading) {
                                Text(recording.recordingDateTime)
                                    .font(.headline)
                                    .scaledToFit()
                            }
                            Spacer()
                            Button(action: {
                                Task {
                                    do {
                                        recordingData = try await getRecording(recordingId: recording.id)
                                        print(recordingData)
                                        path.append(.recording)
                                    } catch RecordingError.invalidURL {
                                        print("RecordingDataError: invalid URL")
                                    } catch RecordingError.serverError {
                                        print("RecordingDataError: internal server error")
                                    } catch RecordingError.invalidData {
                                        print("RecordingDataError: server returned invalid data")
                                    } catch {
                                        print("RecordingDataError: An unexpected error occurred")
                                    }
                                }
                            }) {
                                Image(systemName: "chevron.right")
                                    .foregroundColor(.black)
                            }
                            .frame(width: 34)
                        }
                    }
                    .listStyle(PlainListStyle())
                    .refreshable{
                        listRecordings()
                    }
                }
            }
        }
        .sheet(isPresented: $btPopUp) {
            BluetoothView(bluetoothManager: bluetoothManager)
        }
        .preferredColorScheme(.light)
        .onAppear(perform: listRecordings)
        .onChange(of: isLoggedIn) {
            listRecordings()
        }
    }

    func listRecordings() {
        guard let url = URL(string: "https://heartlink-652851748566.northamerica-northeast2.run.app/ListRecordingsApp?patientid=\(patient.patientId)") else {
            print("invalid URL")
            return
        }
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        URLSession.shared.dataTask(with: request) { data, response, error in
            guard error == nil else {
                print("list recordings: error listing recordings")
                return
            }
            guard let data = data else {
                print("list recordings: did not receive data")
                return
            }
            guard let response = response as? HTTPURLResponse, response.statusCode == 200 else {
                print("list recordings: request failed:")
                return
            }
            do {
                let decoder = JSONDecoder()
                decoder.keyDecodingStrategy = .convertFromSnakeCase
                recordingsList = try decoder.decode(RecordingsList.self, from: data)
            } catch {
                print("list recordings: error decoding")
            }
        }.resume()
    }
}

#Preview {
    MainNavView().environment(\.colorScheme, .light)
}
