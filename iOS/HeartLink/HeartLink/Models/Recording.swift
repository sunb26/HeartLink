//
//  Recording.swift
//  HeartLink
//
//  Created by Ben Sun on 2024-11-10.
//

import Foundation

// Define the Data Model
struct RecordingsList: Codable {
    var widgets: [RecordingWidget]
}

struct RecordingWidget: Identifiable, Codable {
    let id: UInt64
    let recordingDateTime: String
}

struct RecordingData: Codable {
    let recordingId: UInt64
    var status: String
    let physicianComments: String
    let downloadUrl: String
}

struct RecordingSubmission: Codable {
    let recordingId: UInt64
}

enum RecordingError: Error {
    case invalidURL
    case invalidData
    case recordNotFound
    case serverError
}
