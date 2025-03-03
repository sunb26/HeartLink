// firebasedb.go file
// code loosely derived from: https://medium.com/@vubon.roy/lets-integrate-the-firebase-realtime-database-with-golang-7c065a7b7313

package firebasedb

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type FireDB struct {
	*db.Client
}

var fireDB FireDB // define struct instance globally

// connect to firebase database
func (db *FireDB) Connect() error {

	// Get the current file path
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("Unable to get current file info")
	}

	// Get the directory of the current file
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))

	errEnv := godotenv.Load(filepath.Join(rootDir, ".env")) // load environment variables from .env file
	if errEnv != nil {
		log.Fatalf("Error loading .env file: %v", errEnv)
	}
	fmt.Print("GOOGLE_APPLICATION_CREDENTIALS: ", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"), "\n") // TESTING
	fmt.Print("FIREBASE_DATABASE_URL: ", os.Getenv("FIREBASE_DATABASE_URL"), "\n")                   // TESTING

	ctx := context.Background()
	opt := option.WithCredentialsFile(rootDir + os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	config := &firebase.Config{DatabaseURL: os.Getenv("FIREBASE_DATABASE_URL")}
	app, err := firebase.NewApp(ctx, config, opt)
	// app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		return fmt.Errorf("error initializing database: %v", err)
	}

	db.Client = client

	return nil

}

func FirebaseDB() *FireDB {
	return &fireDB
}
