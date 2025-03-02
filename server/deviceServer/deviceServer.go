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

// SEE IF BETTER WAY THAN REDEFINING THIS AGAIN
// const keyServerAddr = "serverAddr" // const string used as key

// GetEndpoint1 function (WILL DELETE)
func GetEndpoint1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Start heartlinkServer endpoint1\n\n") // Write start message to response writer "w"

	// ctx := r.Context() // r.Context used to access the context of the request

	// r.URL.Query used to access the query parameters in the URL
	hasFirst := r.URL.Query().Has("first") // the Has method checks if a query parameter exists (returns Bool)
	first := r.URL.Query().Get("first")    // the Get method retrieves the value of a query parameter (returns String)
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	// Printf writes the arguments extracted from URL to the console
	fmt.Printf("getEndpoint1 arguments: first(%t)=%s, second(%t)=%s\n",
		// ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second)

	// Fprintf writes the arguments extracted from URL to the response writer "w"
	fmt.Fprintf(w, "getEndpoint1 arguments: first(%t)=%s, second(%t)=%s\n",
		// ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second)

	fmt.Fprint(w, "End heartlinkServer endpoint1\n") // Write end message to response writer "w"
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
