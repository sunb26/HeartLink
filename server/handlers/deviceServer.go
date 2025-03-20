package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"heartlinkServer/firebasedb"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// type requestBody struct { // CAN ADD "STRUCT TAGS" TO THIS IF WILL HELP WITH DECODING STRUCT
// 	UserID       string
// 	WavFileBytes uint64 // byte array -> a slice of uint8 values
// }

// POSTRawAudioFile function (OLD)
/*
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

	// No parameters in POST endpoint to read

	// --- extracting request body fields --- //

	userID := req.UserID             // extract userID from request body
	WavFileBytes := req.WavFileBytes // extract wavFileBytes from request body
	//json.Unmarshal([]uint64(req.WavFileBytes), &req.WavFileBytes)

	fmt.Fprint(w, "POSTRawAudioFile request body:\n") // Write request body to response writer "w"
	// fmt.Fprint(w, "userID: ", userID, "\n")           // Write userID to response writer "w"
	fmt.Printf("userID: %v\n", userID) // Write userID to response writer "w"
	fmt.Printf("wavFileBytes: %v\n", WavFileBytes)

	// fmt.Fprint(w, "wavFileBytes: ", req.WavFileBytes[0], "\n")   // Write wavFileBytes to response writer "w"
	//fmt.Fprint(w, "wavFileBytes: ", req.WavFileBytes[1], "\n\n") // second byte in wavFileBytes
	// fmt.Fprint(w, "wavFileBytes: ", string(wavFileBytes), "\n\n") // Write wavFileBytes to response writer "w"

	// --- POST request (from ESP32 to server) requirements --- //



	// --- establishing firebase connection --- //

	errFB := firebasedb.FirebaseDB().Connect() // connect to firebase storage
	if errFB != nil {
		log.Println(errFB)
		return
	}

	// TEST firebase connection via "get all" call
	// testErr := firebasedb.FirebaseDB().GetAllFilesFirebase()
	// if testErr != nil {
	// 	log.Println(testErr)
	// }

	// r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB



	fmt.Fprint(w, "POSTRawAudioFile Endpoint - End\n")
}
*/

func POSTRawAudioFile(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("\nReceiveMultipartForm Endpoint - Start\n\n") // TESTING

	// ensure receiving POST request
	if r.Method != "POST" {
		log.Println("invalid http request type - should be POST request - instead is", r.Method)
	}

	fmt.Fprint(w, "POSTRawAudioFile Endpoint - Start\n") // Arduino expects response when sends POST request

	err := r.ParseMultipartForm(32 << 15) // set 1 MB max on input file size
	if err != nil {
		log.Printf("Error parsing multipart form: %v\n", err)
	}

	// define struct to send response to client
	type uploadedFile struct {
		Size        int64  `json:"size"`
		ContentType string `json:"content_type"`
		Filename    string `json:"filename"`
		FileContent string `json:"file_content"`
	}

	var newFile uploadedFile

	// connect to firebase storage (do this before receive file)
	errFB := firebasedb.FirebaseDB().Connect()
	if errFB != nil {
		log.Printf("Error connecting to Firebase storage %v\n", errFB)
		return
	}

	for _, fheaders := range r.MultipartForm.File {

		for _, headers := range fheaders {

			file, err := headers.Open()
			if err != nil {
				log.Printf("Error opening file: %v\n", err)
				return
			}

			defer file.Close()

			buff := make([]byte, 512)

			file.Read(buff)

			file.Seek(0, 0)

			newFile.ContentType = http.DetectContentType(buff)

			var sizeBuff bytes.Buffer
			fileSize, err := sizeBuff.ReadFrom(file)
			if err != nil {
				log.Printf("Error reading file: %v\n", err)
				return
			}

			fmt.Fprint(w, "server: file size: ", fileSize, "\n") // TESTING - proves endpoint receiving entire file

			file.Seek(0, 0)

			newFile.Size = fileSize

			key := uuid.New() // generate key to act as access token in firebase storage

			newFile.Filename = key.String() + "_" + headers.Filename

			contentBuf := bytes.NewBuffer(nil)

			if _, err := io.Copy(contentBuf, file); err != nil {
				log.Printf("Error copying file: %v\n", err)
				return
			}

			newFile.FileContent = contentBuf.String()

			filePath := "recordings/" + newFile.Filename
			publicURL, err := firebasedb.FirebaseDB().UploadWAVToFirebase([]byte(newFile.FileContent), filePath) // upload wav file to firebase storage
			if err != nil {
				log.Printf("Error uploading WAV file to Firebase Storage: %v\n", err)
			}

			fmt.Printf("Successfully uploaded file to Firebase Storage at: %s\n", publicURL) // TESTING
			fmt.Println()

		}

	}

	userIDWav := newFile.Filename
	userID := strings.Split(strings.Split(userIDWav, "_")[2], ".")[0] // extract userID from filename
	fmt.Printf("userID: %v\n", userID)                                // TESTING - will need to upload userID to relational database

	// populate data to send back to client
	data := make(map[string]interface{})
	data["form_field_value"] = newFile.Filename
	data["status"] = 200
	// data["file_stats"] = newFile // uncomment this if want to send client all file stats

	// send data back to the client
	if err = json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON: %v\n", err)
		return
	}

	fmt.Fprint(w, "POSTRawAudioFile Endpoint - End\n")    // TESTING
	fmt.Printf("\nReceiveMultipartForm Endpoint - End\n") // TESTING
}
