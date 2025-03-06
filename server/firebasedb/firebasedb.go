// firebasedb.go file

package firebasedb

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"firebase.google.com/go/db"

	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type FireDB struct {
	*db.Client
}

var fireDB FireDB
var bucket *storage.BucketHandle

// UploadWAVToFirebase function
func (db *FireDB) UploadWAVToFirebase(localFilePath, firebaseStoragePath string) error {

	client := db.Client
	fmt.Print("Client: ", client, "\n") // TESTING

	fmt.Printf("Successfully uploaded %s to Firebase Storage at path: %s\n", localFilePath, firebaseStoragePath)
	return nil

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

	fmt.Printf(`{"files": %q}`, files)

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
	opt := option.WithCredentialsFile(rootDir + os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

	client, err := storage.NewClient(ctx, opt)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

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
