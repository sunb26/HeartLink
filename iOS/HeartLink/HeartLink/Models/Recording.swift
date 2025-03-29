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
    let date: String
}

struct RecordingData: Codable {
    let id: UInt64
    let date: String
    var viewStatus: String
    let comments: String
    let fileURL: String
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
