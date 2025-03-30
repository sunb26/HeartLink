//
//  User.swift
//  HeartLink
//
//  Created by Ben Sun on 2025-01-20.
//

import Foundation

struct User: Decodable {
    let patientId: UInt64
}

enum LoginError: Error {
    case invalidURL
    case invalidCredentials
    case invalidResponse
    case invalidData
    case serverError
}
