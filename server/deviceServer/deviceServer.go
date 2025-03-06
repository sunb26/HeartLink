// deviceServer.go file

package deviceServer

import (
	"encoding/json"
	"errors"
	"fmt"
	"heartlinkServer/firebasedb"
	"io"
	"log"
	"net/http"
)

type requestBody struct { // CAN ADD "STRUCT TAGS" TO THIS IF WILL HELP WITH DECODING STRUCT
	UserID       string
	WavFileBytes []byte // byte array -> a slice of uint8 values
}

// POSTRawAudioFile function
func POSTRawAudioFile(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "POSTRawAudioFile Endpoint - Start\n\n")

	r.Body = http.MaxBytesReader(w, r.Body, 1*1024*1024) // limit the size of the request body to 1 MB

	decoder := json.NewDecoder(r.Body) // create a new JSON decoder

	var req requestBody // create a new requestBody struct

	err := decoder.Decode(&req) // decode the JSON request body into the requestBody struct

	if err != nil {
		var syntaxError *json.SyntaxError
		var MaxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)

		case errors.As(err, &MaxBytesError):
			msg := fmt.Sprintf("Request body is too large. Max size is %d bytes", MaxBytesError.Limit)
			http.Error(w, msg, http.StatusRequestEntityTooLarge)

		default:
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = decoder.Decode(&struct{}{}) // check if there are any extra JSON objects/characters in the request body
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain one JSON object"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	/* No parameters in POST endpoint to read */

	// --- extracting request body fields --- //

	userID := req.UserID // extract userID from request body
	// wavFileBytes := req.WavFileBytes // extract wavFileBytes from request body
	json.Unmarshal([]byte(req.WavFileBytes), &req.WavFileBytes)

	fmt.Fprint(w, "POSTRawAudioFile request body:\n")            // Write request body to response writer "w"
	fmt.Fprint(w, "userID: ", userID, "\n")                      // Write userID to response writer "w"
	fmt.Fprint(w, "wavFileBytes: ", req.WavFileBytes[0], "\n")   // Write wavFileBytes to response writer "w"
	fmt.Fprint(w, "wavFileBytes: ", req.WavFileBytes[1], "\n\n") // second byte in wavFileBytes
	// fmt.Fprint(w, "wavFileBytes: ", string(wavFileBytes), "\n\n") // Write wavFileBytes to response writer "w"

	// --- POST request (from ESP32 to server) requirements --- //

	/*
		1. no arguments in URL to extract
		2. request body in JSON format with the following fields:
			a. userID (string) - ESP32 should get this from app (via BT)
			b. wavFileBytes (byte array) - however .wav file is broken down into bytes, send here
		3. ESP32 will receive response from server (e.g. success code) to know worked
	*/

	// --- establishing firebase connection --- //

	errFB := firebasedb.FirebaseDB().Connect() // connect to firebase storage
	if errFB != nil {
		log.Println(errFB)
		return
	}

	// TEST firebase connection via "get all" call
	testErr := firebasedb.FirebaseDB().GetAllFilesFirebase()
	if testErr != nil {
		log.Println(testErr)
	}

	// testWavFilepath := "testWavFile1.wav"                                                          // path to test wav file
	// output := firebasedb.FirebaseDB().UploadWAVToFirebase(testWavFilepath, "recordings/test1.wav") // upload wav file to firebase storage
	// if output != nil {
	// 	log.Fatalf("Error uploading WAV file to Firebase Storage: %v", output)
	// }

	fmt.Fprint(w, "POSTRawAudioFile Endpoint - End\n")
}
