//
//  ContentView.swift
//  HeartLink
//
//  Created by Ben Sun on 2024-10-14.
//

import SwiftUI

struct LoginView: View {
    @Binding var path: [PageActions]
    @Binding var patient: User
    @State private var userId: String = ""
    @State private var password: String = ""
    @State private var isLoading: Bool = false
    @State private var loginFailed: Bool = false

    var body: some View {
        ZStack {
            Color.white
                .ignoresSafeArea() // Ensure full white background

            VStack {
                HStack {
                    Image("logo")
                        .resizable()
                        .aspectRatio(contentMode: .fit)
                        .frame(width: 250, height: 250)
                    // #d32f2f
                }
                .padding(.bottom, 50)

                VStack {
                    TextField("Username", text: $userId)
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

                    if loginFailed {
                        Text("Invalid Credentials")
                            .padding(.top, 10)
                    }

                    Button(action: {
                        loginFailed = false
                        if userId.isEmpty || password.isEmpty {
                            loginFailed = true
                        } else {
                            Task {
                                await performLogin()
                            }
                        }
                    }) {
                        if isLoading {
                            ProgressView()
                        } else {
                            Text("Login")
                                .padding()
                                .frame(maxWidth: .infinity)
                                .background(Color("custom-red"))
                                .foregroundColor(.white)
                                .cornerRadius(8)
                                .shadow(color: .gray, radius: 5, x: 0, y: 2)
                        }
                    }
                    .padding(.top, 60)
                    .disabled(isLoading)
                }
                .padding(.horizontal, 15)
            }
            .padding()
            .ignoresSafeArea(edges: .all)
        }
    }

    private func performLogin() async {
        isLoading = true

        do {
            patient = try await getUser(username: userId, password: password)
            isLoading = false
            path.removeLast(path.count)
        } catch LoginError.invalidURL {
            print("LoginError: invalid URL")
        } catch LoginError.invalidResponse {
            print("LoginError: response parsing error")
        } catch LoginError.invalidCredentials {
            print("LoginError: invalid username or password")
            loginFailed = true
        } catch LoginError.serverError {
            print("LoginError: internal server error")
        } catch LoginError.invalidData {
            print("LoginError: server returned invalid data")
        } catch {
            print("LoginError: An unexpected error occurred")
        }

        isLoading = false
    }
}

struct ResetPasswordView: View {
    var body: some View {
        Text("Reset your password here")
    }
}

#Preview {
    @Previewable @State var path: [PageActions] = [.login]
    @Previewable @State var patient: User = User(patientId: 1)

    LoginView(path: $path, patient: $patient)
}
