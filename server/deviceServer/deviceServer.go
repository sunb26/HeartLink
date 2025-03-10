// deviceServer.go file

package deviceServer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"heartlinkServer/firebasedb"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type requestBody struct { // CAN ADD "STRUCT TAGS" TO THIS IF WILL HELP WITH DECODING STRUCT
	UserID       string
	WavFileBytes uint64 // byte array -> a slice of uint8 values
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
			fmt.Printf("default error: %v\n", err)
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

	userID := req.UserID             // extract userID from request body
	WavFileBytes := req.WavFileBytes // extract wavFileBytes from request body
	/* json.Unmarshal([]byte(req.WavFileBytes), &req.WavFileBytes) */
	//json.Unmarshal([]uint64(req.WavFileBytes), &req.WavFileBytes)

	fmt.Fprint(w, "POSTRawAudioFile request body:\n") // Write request body to response writer "w"
	// fmt.Fprint(w, "userID: ", userID, "\n")           // Write userID to response writer "w"
	fmt.Printf("userID: %v\n", userID) // Write userID to response writer "w"
	fmt.Printf("wavFileBytes: %v\n", WavFileBytes)

	/* fmt.Fprint(w, "wavFileBytes: ", req.WavFileBytes, "\n") // Write wavFileBytes to response writer "w" */
	// fmt.Fprint(w, "wavFileBytes: ", req.WavFileBytes[0], "\n")   // Write wavFileBytes to response writer "w"
	//fmt.Fprint(w, "wavFileBytes: ", req.WavFileBytes[1], "\n\n") // second byte in wavFileBytes
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

	// r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB

	/*
		testWavFilepath := "testWavFile1.wav"                                                       // path to test wav file (with respect to current file)
		publicURL, err := firebasedb.FirebaseDB().UploadWAVToFirebase()(testWavFilepath, SOMETHING) // upload wav file to firebase storage
		if err != nil {
			log.Fatalf("Error uploading WAV file to Firebase Storage: %v", err)
		}
		fmt.Fprintf(w, `{"status": "success", "url": %q}`, publicURL)
	*/

	fmt.Fprint(w, "POSTRawAudioFile Endpoint - End\n")
}

//
//
//
//
//
//

func ReceiveMultipartForm(w http.ResponseWriter, r *http.Request) {

	// ensure receiving POST request
	if r.Method != "POST" {
		log.Fatalln(errors.New("invalid method"))
		return
	}

	fmt.Fprint(w, "POSTRawAudioFile Endpoint - Start\n")    // NEED THIS (Arduino expects resposne when sends POST request)
	fmt.Printf("ReceiveMultipartForm Endpoint - Start\n\n") // testing

	err := r.ParseMultipartForm(32 << 15) // 1 MB max
	if err != nil {
		fmt.Printf("Error parsing multipart form: %v\n", err)
		log.Fatal(err) // log better error message
	}

	// define struct to send back response to the client
	type uploadedFile struct {
		Size        int64  `json:"size"`
		ContentType string `json:"content_type"`
		Filename    string `json:"filename"`
		FileContent string `json:"file_content"`
	}

	var newFile uploadedFile

	// create uploads directory (if it doesn't exist) - TESTING ONLY (will not be saving locally in the end)
	if err := os.MkdirAll("uploads", 0755); err != nil {
		fmt.Printf("Error creating uploads directory: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, fheaders := range r.MultipartForm.File {

		for _, headers := range fheaders {

			file, err := headers.Open()

			if err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				fmt.Printf("Error opening file: %v", err)
				return
			}

			defer file.Close()

			// detect contentType

			buff := make([]byte, 512)

			file.Read(buff)

			file.Seek(0, 0)

			contentType := http.DetectContentType(buff)

			newFile.ContentType = contentType

			// get file size

			var sizeBuff bytes.Buffer
			fileSize, err := sizeBuff.ReadFrom(file)
			if err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				fmt.Printf("Error reading file: %v", err)
				return
			}

			fmt.Fprint(w, "server: file size: ", fileSize, "\n") // proves endpoint receiving entire file

			file.Seek(0, 0) // reset read/write back to position 0

			newFile.Size = fileSize // write how large file is to newFile

			newFile.Filename = headers.Filename

			contentBuf := bytes.NewBuffer(nil)

			if _, err := io.Copy(contentBuf, file); err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				fmt.Printf("Error copying file: %v", err)
				return
			}

			newFile.FileContent = contentBuf.String()

			// upload file to LOCAL storage (for testing purposes)
			// headers.Filename = "testFileManual2.wav" // TESTING
			filePath := filepath.Join("uploads", headers.Filename)
			err = os.WriteFile(filePath, []byte(newFile.FileContent), 0644)
			if err != nil {
				fmt.Printf("Error saving file: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Printf("Successfully saved file: %s\n", filePath) // TESTING

		}

	}

	userIDWav := newFile.Filename
	userID := strings.Split(strings.Split(userIDWav, "_")[1], ".")[0] // final userID value
	fmt.Printf("userID: %v\n", userID)                                // testing

	fmt.Fprint(w, "server: about to send data\n") // TESTING

	data := make(map[string]interface{})

	// define the data to send back to the client
	data["form_field_value"] = newFile.Filename
	data["status"] = 200
	// data["file_stats"] = newFile // UNCOMMENT THIS (if want to send client all file stats)

	// send data back to the client
	if err = json.NewEncoder(w).Encode(data); err != nil {
		fmt.Printf("Error encoding JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "server: data sent (2)\n") // TESTING

	fmt.Fprint(w, "POSTRawAudioFile Endpoint - End\n")    // TESTING (not needed)
	fmt.Printf("ReceiveMultipartForm Endpoint - End\n\n") // TESTING (not needed)
}
