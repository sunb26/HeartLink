////
////  RecordingActions.swift
////  HeartLink
////
////  Created by Ben Sun on 2025-01-23.
////
//
import Foundation

func getRecording(recordingId: UInt64) async throws -> RecordingData {
    guard let url = URL(string: "https://heartlink-652851748566.northamerica-northeast2.run.app/LoadRecordingInfoApp?recordingid=\(recordingId)") else {
        throw RecordingError.invalidURL
    }

    let (data, response) = try await URLSession.shared.data(from: url)

    guard let response = response as? HTTPURLResponse else {
        throw RecordingError.serverError
    }

    guard response.statusCode == 200 else {
        if response.statusCode == 404 {
            throw RecordingError.recordNotFound
        } else {
            throw RecordingError.serverError
        }
    }

    do {
        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        return try decoder.decode(RecordingData.self, from: data)
    } catch {
        throw RecordingError.invalidData
    }
}

func submit(submission: RecordingSubmission) async throws {
    guard let url = URL(string: "https://heartlink-652851748566.northamerica-northeast2.run.app/SaveRunAlgorithm") else {
        throw RecordingError.invalidURL
    }
    var request = URLRequest(url: url)
    request.httpMethod = "POST"
    request.addValue("application/json", forHTTPHeaderField: "Content-Type")

    do {
        let jsonData = try JSONEncoder().encode(submission)
        request.httpBody = jsonData
        
        let (_, response) = try await URLSession.shared.data(for: request)

        guard let response = response as? HTTPURLResponse else {
            print("Invalid URL")
            throw RecordingError.serverError
        }

        guard response.statusCode == 200 else {
            print("Response Code: \(response.statusCode)")
            throw RecordingError.serverError
        }
        print("submitted successfully")
        return
    } catch {
        print("failed to encode submission: \(error)")
    }
}

func delete(recordingId: UInt64) async throws {
    guard let url = URL(string: "https://heartlink-652851748566.northamerica-northeast2.run.app/DeleteRecording?recordingid=\(recordingId)") else {
        print("delete recording: invalid URL")
        return
    }
    var request = URLRequest(url: url)
    request.httpMethod = "DELETE"
    do {
        let (_, response) = try await URLSession.shared.data(for: request)
        guard let response = response as? HTTPURLResponse else {
            throw RecordingError.serverError
        }

        guard response.statusCode == 200 else {
            print("failed to delete recording: return code \(response.statusCode)")
            throw RecordingError.serverError
        }
        print("deleted successfully")
        return
    } catch {
        print("error deleting recording")
    }
}
