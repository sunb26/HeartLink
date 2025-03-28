package handlers

import (
	"bytes"
	"fmt"
	"heartlinkServer/firebasedb"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (env *Env) UploadFilterRecording(w http.ResponseWriter, r *http.Request) {

	// ensure receiving POST request
	if r.Method != "POST" {
		log.Println("invalid http request type - should be POST request - instead is", r.Method)
	}

	fmt.Fprint(w, "UploadFilterRecording Endpoint - Start\n") // Arduino expects response when sends POST request

	// set 1 MB max on input file size
	err := r.ParseMultipartForm(32 << 15)
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

	/* COMMENT OUT WHEN COMMITTING - FOR RUNNING WITH LOCAL HOST ONLY */
	// _, currentFile, _, ok := runtime.Caller(0) // get current file path
	// if !ok {
	// 	log.Fatalf("Unable to get current file info")
	// }
	// rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile))) // get root directory of current file
	// err = godotenv.Load(filepath.Join(rootDir, ".env"))              // load environment variables from .env file
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }
	/* END COMMENT OUT WHEN COMMITTING - FOR RUNNING WITH LOCAL HOST ONLY */

	var newFile uploadedFile
	var publicURL string

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

			// generate a POST request to send the unfiltered audio file to the python server for DSP
			var requestBody bytes.Buffer

			multipartWriter := multipart.NewWriter(&requestBody)

			fileWriter, err := multipartWriter.CreateFormFile("audio", newFile.Filename)
			if err != nil {
				log.Printf("Error creating form file: %v\n", err)
			}

			fileContentReader := strings.NewReader(newFile.FileContent)
			_, err = io.Copy(fileWriter, fileContentReader)
			if err != nil {
				log.Printf("Error copying file: %v\n", err)
			}

			multipartWriter.Close()

			req, err := http.NewRequest("POST", os.Getenv("PYTHON_SERVER_URL"), &requestBody)
			if err != nil {
				log.Printf("Error creating POST request: %v\n", err)
			}
			req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Erorr sending POST request to python server: %v\n", err)
			}
			defer resp.Body.Close()

			// upload the filtered audio recording to Firebase
			filePath := "recordings/" + newFile.Filename
			publicURL, err = firebasedb.FirebaseDB().UploadWAVToFirebase(resp.Body, filePath)
			if err != nil {
				log.Printf("Error uploading WAV file to Firebase Storage: %v\n", err)
			}

			fmt.Printf("Successfully uploaded file to Firebase Storage at: %s\n", publicURL) // TESTING
			fmt.Println()

		}

	}

	userIDWav := newFile.Filename
	userID := strings.Split(strings.Split(userIDWav, "_")[2], ".")[0] // extract userID from filename

	// send userID + public firebase recording link to relational database
	tx, err := env.DB.Beginx() // setup database connection
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var status string = "notSubmitted" // always set to notSubmitted when first uploaded
	_, err = tx.Exec("INSERT INTO recordings (patient_id, download_url, recording_datetime, status) VALUES ($1, $2, $3, $4)", userID, publicURL, time.Now(), status)
	if err != nil {
		log.Printf("Error inserting new recording into database: %v\n", err)
	}

	// commit transaction to database
	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// w.WriteHeader(http.StatusOK)

	// populate data to send back to client - not sure if needed/following this convention
	// data := make(map[string]interface{})
	// data["form_field_value"] = newFile.Filename
	// data["status"] = 200

	// send data back to the client
	// if err = json.NewEncoder(w).Encode(data); err != nil {
	// 	log.Printf("Error encoding JSON: %v\n", err)
	// 	return
	// }

	fmt.Fprint(w, "UploadFilterRecording Endpoint - End\n") // TESTING
}
