// firebasedb.go file

package firebasedb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"cloud.google.com/go/storage"
	"firebase.google.com/go/db"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type FireDB struct {
	*db.Client
}

type configFile struct {
	Type                        string
	Project_id                  string
	Private_key_id              string
	Private_key                 string
	Client_email                string
	Client_id                   string
	Auth_uri                    string
	Token_uri                   string
	Auth_provider_x509_cert_url string
	Client_x509_cert_url        string
	Universe_domain             string
}

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

	// get current file path
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("Unable to get current file info")
	}

	// get root directory of current file (based on the current file strcuture)
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))

	// load environment variables from .env file
	err := godotenv.Load(filepath.Join(rootDir, ".env"))
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ctx := context.Background()
	// opt := option.WithCredentialsFile(rootDir + os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

	configFile1 := configFile{
		os.Getenv("Type"),
		os.Getenv("Project_id"),
		os.Getenv("Private_key_id"),
		os.Getenv("Private_key"),
		os.Getenv("Client_email"),
		os.Getenv("Client_id"),
		os.Getenv("Auth_uri"),
		os.Getenv("Token_uri"),
		os.Getenv("Auth_provider_x509_cert_url"),
		os.Getenv("Client_x509_cert_url"),
		os.Getenv("Universe_domain"),
	}

	configFile1JSON, err := json.Marshal(configFile1)
	if err != nil {
		fmt.Printf("Error creating JSON file: %v\n", err)
	}

	opt := option.WithCredentialsJSON(configFile1JSON)

	client, err := storage.NewClient(ctx, opt)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// token, err := client.CustomToken(ctx, os.Getenv("FIREBASE_UID"))

	bucket = client.Bucket(os.Getenv("FIREBASE_STORAGE_BUCKET"))

	// check if bucket exists + created successfully
	_, err = bucket.Attrs(ctx)
	if err != nil {
		log.Fatalf("Failed to get bucket attributes: %v", err)
	}

	fmt.Printf("Successfully connected to Firebase Storage\n") // TESTING

	return nil

}

func FirebaseDB() *FireDB {
	return &fireDB
}
