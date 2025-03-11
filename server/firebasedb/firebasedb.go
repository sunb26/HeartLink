// firebasedb.go file

package firebasedb

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"firebase.google.com/go/db"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type FireDB struct {
	*db.Client
}

// type configFile struct {
// 	Type                        string `json:"type"`
// 	project_id                  string
// 	private_key_id              string
// 	private_key                 string
// 	client_email                string
// 	client_id                   string
// 	auth_uri                    string
// 	token_uri                   string
// 	auth_provider_x509_cert_url string
// 	client_x509_cert_url        string
// 	universe_domain             string
// }

var fireDB FireDB
var bucket *storage.BucketHandle

// UploadWAVToFirebase function
// func (db *FireDB) UploadWAVToFirebase(file multipart.File, storagePath string) (string, error) {
func (db *FireDB) UploadWAVToFirebase(fileContent []byte, storagePath string) (string, error) { // need to figure out the arguments for this function

	ctx := context.Background()

	object := bucket.Object(storagePath)

	// set up writing object to write .wav files
	writer := object.NewWriter(ctx)
	writer.ContentType = "audio/wav"
	// metadata: { firebaseStorageDownloadTokens: uuidv3() }
	// writer.Metadata = map[string]string{
	// 	metadata: {
	// 		firebaseStorageDownloadTokens: uuidv3(),
	// 	},
	// }

	// read wav file data (temporary)
	/* f, err := os.Open(file)
	if err != nil {
		fmt.Printf("os.Open error: %v", err)
		return "", err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Printf("io.ReadAll error: %v", err)
		return "", err
	}*/

	// copy file data to storage
	_, err := io.Copy(writer, bytes.NewReader(fileContent))
	if err != nil {
		fmt.Printf("io.Copy error: %v", err)
		return "", err
	}

	// close writer to finalize upload
	err = writer.Close()
	if err != nil {
		fmt.Printf("Writer.Close error: %v", err)
		return "", err
	}

	// set public access
	err = object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		fmt.Printf("ACL.Set error: %v", err)
		return "", err
	}

	// get public URL
	attrs, err := object.Attrs(ctx)
	if err != nil {
		fmt.Printf("Attrs error: %v", err)
		return "", err
	}

	return attrs.MediaLink, nil

}

// TEST function (will use in other endpoints)
func (db *FireDB) GetAllFilesFirebase() error {

	directory := "recordings/" // hard coded (based on Firebase Storage structure)

	ctx := context.Background()

	query := &storage.Query{
		Prefix:    directory,
		Delimiter: "/",
	}

	var files []string
	it := bucket.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("Error iterating through bucket: %v", err)
			return nil
		}
		files = append(files, attrs.Name)
	}

	fmt.Printf(`{"files": %q`, files)
	fmt.Println()

	return nil

}

// connect to firebase database
func (db *FireDB) Connect() error {

	fmt.Printf("Connecting to Firebase Storage\n") // TESTING

	/* COMMENT OUT WHEN COMMITTING */
	// _, currentFile, _, ok := runtime.Caller(0) // get current file path
	// if !ok {
	// 	log.Fatalf("Unable to get current file info")
	// }
	// rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile))) // get root directory of current file (based on the current file strcuture)
	// err := godotenv.Load(filepath.Join(rootDir, ".env"))             // load environment variables from .env file
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }
	/* END COMMENT OUT WHEN COMMITTING */

	ctx := context.Background()
	// opt := option.WithCredentialsFile(rootDir + os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

	// privateKeyRaw := os.Getenv("private_key")
	// fmt.Println("Raw private key length:", len(privateKeyRaw))
	// fmt.Println("Raw private key first 20 chars:", privateKeyRaw[:40])

	// configFile1 := configFile{
	// 	os.Getenv("type"),
	// 	os.Getenv("project_id"),
	// 	os.Getenv("private_key_id"),
	// 	strings.ReplaceAll(os.Getenv("private_key"), "\\n", "\n"),
	// 	os.Getenv("client_email"),
	// 	os.Getenv("client_id"),
	// 	os.Getenv("auth_uri"),
	// 	os.Getenv("token_uri"),
	// 	os.Getenv("auth_provider_x509_cert_url"),
	// 	os.Getenv("client_x509_cert_url"),
	// 	os.Getenv("universe_domain"),
	// }

	// configPrivateKey := configFile1.private_key
	// fmt.Println("After: private key length:", len(configPrivateKey))
	// fmt.Println("After: private key chars:", configPrivateKey)

	// fmt.Printf("json file: %v\n", configFile1.Type)

	// configFile1JSON, err := json.Marshal(configFile1)
	// if err != nil {
	// 	fmt.Printf("Error creating JSON file: %v\n", err)
	// }

	// fmt.Printf("json file: %v\n", configFile1JSON)

	// opt := option.WithCredentialsJSON(configFile1JSON)

	// tempFile, err := os.CreateTemp("", "google-credentials-*.json")
	// if err != nil {
	// 	log.Fatalf("Failed to create temp file: %v", err)
	// }
	// defer os.Remove(tempFile.Name())

	// if _, err := tempFile.Write(configFile1JSON); err != nil {
	// 	log.Fatalf("Failed to write to temp file: %v", err)
	// }
	// if err := tempFile.Close(); err != nil {
	// 	log.Fatalf("Failed to close temp file: %v", err)
	// }

	// opt := option.WithCredentialsFile(tempFile.Name())

	//try

	credentialJSON := []byte(fmt.Sprintf(`{
		"type": %q,
		"project_id": %q,
		"private_key_id": %q,
		"private_key": %q,
		"client_email": %q,
		"client_id": %q,
		"auth_uri": %q,
		"token_uri": %q,
		"auth_provider_x509_cert_url": %q,
		"client_x509_cert_url": %q,
		"universe_domain": %q
	}`,
		os.Getenv("type"),
		os.Getenv("project_id"),
		os.Getenv("private_key_id"),
		strings.ReplaceAll(os.Getenv("private_key"), "\\n", "\n"),
		os.Getenv("client_email"),
		os.Getenv("client_id"),
		os.Getenv("auth_uri"),
		os.Getenv("token_uri"),
		os.Getenv("auth_provider_x509_cert_url"),
		os.Getenv("client_x509_cert_url"),
		os.Getenv("universe_domain"),
	))

	opt := option.WithCredentialsJSON(credentialJSON)

	client, err := storage.NewClient(ctx, opt)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// token, err := client.CustomToken(ctx, os.Getenv("FIREBASE_UID"))

	bucket = client.Bucket(os.Getenv("FIREBASE_STORAGE_BUCKET"))

	// check if bucket exists + created successfully
	_, err = bucket.Attrs(ctx)
	if err != nil {
		fmt.Printf("Failed to get bucket attributes - error message: %v\n", err)
		// log.Fatalf("Failed to get bucket attributes: %v", err)
	}

	fmt.Printf("Successfully connected to Firebase Storage\n") // TESTING

	return nil

}

func FirebaseDB() *FireDB {
	return &fireDB
}
