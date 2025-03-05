//
//  RecLogView.swift
//  HeartLink
//
//  Created by Ben Sun on 2024-12-20.
//
// Displays the data for a specific recording.

import AVFoundation
import FirebaseStorage
import SwiftUI

struct RecLogView: View {
    @Binding var path: [PageActions]
    @Binding var recording: RecordingData
    @State var isPlaying: Bool = false
    @State private var player = AVPlayer()
    @State private var duration: Double = 0
    @State private var currentTime: Double = 0
    @State private var timer: Timer?
    @State private var confirmSubmit: Bool = false
    @State private var confirmDelete: Bool = false

    var body: some View {
        VStack {
            VStack {
                VStack(alignment: .leading) {
                    Text("Recording")
                        .font(.title).fontWeight(.bold)
                        .foregroundColor(.black)
                    Text("ID: \(recording.id)")
                    if recording.viewStatus == "viewed" {
                        VStack(alignment: .leading) {
                            HStack {
                                Text("View Status:")
                                Text("Viewed")
                                    .foregroundStyle(.green)
                                    .bold()
                            }
                            VStack(alignment: .leading) {
                                Text("Comments: ")
                                    .font(.title2)
                                    .bold()
                                ScrollView {
                                    Text(recording.comments)
                                        .multilineTextAlignment(.leading)
                                        .font(.system(size: 20))
                                }
                            }
                            .padding(.top, 25)
                        }
                    } else if recording.viewStatus == "notSubmitted" {
                        VStack(alignment: .leading) {
                            HStack {
                                Text("View Status:")
                                Text("Not Submitted")
                                    .foregroundStyle(.red)
                                    .bold()
                            }
                            .padding(.bottom, 20)

                            HStack {
                                Button("Submit", action: {
                                    self.confirmSubmit = true
                                })
                                .buttonStyle(.bordered)
                                .confirmationDialog("Are you sure you want to submit?",
                                                    isPresented: $confirmSubmit,
                                                    titleVisibility: .visible) {
                                    Button("Submit", action: {
                                        Task {
                                            do {
                                                let submission = RecordingSubmission(recordingId: recording.id, url: recording.fileURL)
                                                try await submit(submission: submission)
                                                recording.viewStatus = "submitted"
                                            } catch {
                                                print("error submitting recording: \(error)")
                                            }
                                        }
                                    })
                                }
                                Spacer()
                                Button("Delete", role: .destructive, action: {
                                    self.confirmDelete = true
                                })
                                .buttonStyle(.bordered)
                                .confirmationDialog("Are you sure? This action cannot be undone.",
                                                    isPresented: $confirmDelete,
                                                    titleVisibility: .visible) {
                                    Button("Delete Recording", role: .destructive) {
                                        Task {
                                            do {
                                                try await delete(recordingId: recording.id)
                                                path.removeLast(path.count)
                                            } catch {
                                                print("error deleting record: \(error)")
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    } else if recording.viewStatus == "pending" {
                        HStack {
                            Text("View Status:")
                            Text("Submitted - Under Review")
                                .foregroundStyle(.yellow)
                                .bold()
                        }
                    } else {
                        HStack {
                            Text("View Status:")
                            Text("Error")
                                .foregroundStyle(.yellow)
                                .bold()
                        }
                    }
                }
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .topLeading)
            .padding(20)

            // Progress Bar
            if duration > 0 {
                Slider(value: $currentTime, in: 0 ... duration, onEditingChanged: sliderEditingChanged)
                    .padding()

                // Timestamps
                HStack {
                    Text(formatTime(currentTime)) // Current time
                    Spacer()
                    Text(formatTime(duration)) // Total duration
                }
                .padding(.horizontal)
                Button(action: self.playPause, label: {
                    Image(systemName: isPlaying ? "pause.circle.fill" : "play.circle.fill").resizable()
                }).frame(width: 70, height: 70, alignment: .center)
            } else {
                Text("Loading ...")
            }
        }
        .onAppear(perform: setupPlayer)
        .onDisappear(perform: cleanup)
    }

    func setupPlayer() {
        let storage = Storage.storage().reference(forURL: recording.fileURL)
        storage.downloadURL { url, error in
            if error != nil {
                print("display recording data: \(String(describing: error))")
            } else {
                do {
                    try
                        // Playback even with notifications silenced
                        AVAudioSession.sharedInstance().setCategory(AVAudioSession.Category.playback)
                } catch {
                    print("display recording data: \(String(describing: error))")
                }
                self.player = AVPlayer(url: url!)

                Task {
                    let totalTime = try await player.currentItem?.asset.load(.duration)
                    self.duration = totalTime?.seconds ?? 0
                }
            }
        }
    }

    func playPause() {
        isPlaying.toggle()
        if isPlaying {
            player.play()
            startTimer()
        } else {
            player.pause()
            stopTimer()
        }
    }

    func sliderEditingChanged(editingStarted: Bool) {
        if editingStarted {
            stopTimer()
        } else {
            let targetTime = CMTime(seconds: currentTime, preferredTimescale: 600)
            player.seek(to: targetTime)
            if isPlaying {
                startTimer()
            }
        }
    }

    func startTimer() {
        timer = Timer.scheduledTimer(withTimeInterval: 0.1, repeats: true) { _ in
            if let currentTimeValue = player.currentItem?.currentTime().seconds {
                self.currentTime = currentTimeValue
                if self.currentTime == self.duration {
                    self.isPlaying = false
                    self.currentTime = 0
                    self.player.seek(to: .zero)
                    stopTimer()
                }
            }
        }
    }

    func stopTimer() {
        timer?.invalidate()
        timer = nil
    }

    func formatTime(_ time: Double) -> String {
        let minutes = Int(time) / 60
        let seconds = Int(time) % 60
        return String(format: "%02d:%02d", minutes, seconds)
    }

    func cleanup() {
        player.pause()
        stopTimer()
        isPlaying = false
    }
}

#Preview {
    @Previewable @State var path: [PageActions] = [.recording]
    @Previewable @State var record = RecordingData(id: 0, date: "2025-01-01", viewStatus: "notSubmitted", comments: "Test", fileURL: "gs://heartlink-6fee0.firebasestorage.app/recordings/chillin39-20915.wav")
    RecLogView(path: $path, recording: $record)
}
