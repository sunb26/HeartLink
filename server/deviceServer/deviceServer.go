// deviceServer.go file

package deviceServer

import (
	"encoding/json"
	"errors"
	"fmt"
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
	fmt.Fprint(w, "POSTRawAudioFile Endpoint - Start\n\n") // Write start message to response writer "w"

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
	// Fprintf writes the arguments extracted from URL to the response writer "w"
	fmt.Fprintf(w, "POSTRawAudioFile arguments: N/A\n\n")

	userID := req.UserID // extract userID from request body
	// wavFileBytes := req.WavFileBytes // extract wavFileBytes from request body
	json.Unmarshal([]byte(req.WavFileBytes), &req.WavFileBytes)

	fmt.Fprint(w, "POSTRawAudioFile request body:\n")            // Write request body to response writer "w"
	fmt.Fprint(w, "userID: ", userID, "\n")                      // Write userID to response writer "w"
	fmt.Fprint(w, "wavFileBytes: ", req.WavFileBytes[0], "\n")   // Write wavFileBytes to response writer "w"
	fmt.Fprint(w, "wavFileBytes: ", req.WavFileBytes[1], "\n\n") // second byte in wavFileBytes
	// fmt.Fprint(w, "wavFileBytes: ", string(wavFileBytes), "\n\n") // Write wavFileBytes to response writer "w"

	fmt.Fprint(w, "POSTRawAudioFile Endpoint - End\n") // Write end message to response writer "w"
}
