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
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type FireDB struct {
	*db.Client
}

var fireDB FireDB
var bucket *storage.BucketHandle

// UploadWAVToFirebase function
func (db *FireDB) UploadWAVToFirebase(fileContent []byte, storagePath string) (string, error) {

	ctx := context.Background()
	key := uuid.New() // generate key to act as access token in firebase storage

	fmt.Println("key value (in firebasedb.go):", key.String()) // TESTING

	object := bucket.Object(storagePath)

	// set up object to write .wav files
	writer := object.NewWriter(ctx)
	writer.ContentType = "audio/wav"
	writer.ChunkSize = 0
	writer.ObjectAttrs.Metadata = map[string]string{
		"firebaseStorageDownloadTokens": key.String(), // create access token from uniquely generated key
	}

	// copy file data to storage
	_, err := io.Copy(writer, bytes.NewReader(fileContent))
	if err != nil {
		log.Printf("io.Copy error: %v\n", err)
		return "", err
	}

	// close writer to finalize upload
	err = writer.Close()
	if err != nil {
		log.Printf("object.Close error: %v\n", err)
		return "", err
	}

	// set public access
	err = object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		log.Printf("ACL.Set error: %v\n", err)
		return "", err
	}

	// get object attributes to return public URL
	attrs, err := object.Attrs(ctx)
	if err != nil {
		log.Printf("Attrs error: %v\n", err)
		return "", err
	}

	return attrs.MediaLink, nil

}

// TEST function (will use in other endpoints)
/*
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
*/

// connect to firebase database
func (db *FireDB) Connect() error {

	fmt.Printf("Connecting to Firebase Storage\n") // TESTING

	ctx := context.Background()

	/* COMMENT OUT WHEN COMMITTING - FOR RUNNING WITH LOCAL HOST ONLY */
	// _, currentFile, _, ok := runtime.Caller(0) // get current file path
	// if !ok {
	// 	log.Fatalf("Unable to get current file info")
	// }
	// rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile))) // get root directory of current file (based on the current file structure)
	// err := godotenv.Load(filepath.Join(rootDir, ".env"))             // load environment variables from .env file
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }
	/* END COMMENT OUT WHEN COMMITTING - FOR RUNNING WITH LOCAL HOST ONLY */

	// define credentials JSON file to store all environment variables
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
		log.Fatalf("Failed to create client: %v\n", err)
	}

	bucket = client.Bucket(os.Getenv("FIREBASE_STORAGE_BUCKET"))

	// check bucket created successfully
	_, err = bucket.Attrs(ctx)
	if err != nil {
		log.Fatalf("Failed to get bucket attributes - error message: %v\n", err)
	}

	fmt.Printf("Successfully connected to Firebase Storage\n") // TESTING

	return nil

}

func FirebaseDB() *FireDB {
	return &fireDB
}
