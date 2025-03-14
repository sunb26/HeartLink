// main.go file

package main

import (
	"context"
	"errors"
	"fmt"
	"heartlinkServer/handlers"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("Begin main function")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	// connect to the SQL database
	db, err := sqlx.Connect("postgres", os.Getenv("SQL_DSN"))
	if err != nil {
		log.Fatalf("failed to open to database: %v", err)
	}

	// Instantiate dependencies to be injected for the handlers
	env := &handlers.Env{DB: db}

	mux := http.NewServeMux() // create custom multiplexer to handle incoming requests

	// each HandleFunc is used to handle a specific endpoint
	mux.HandleFunc("/POSTRawAudioFile", handlers.POSTRawAudioFile)
	mux.HandleFunc("/createPhysician", env.CreatePhysician)

	ctx := context.Background()

	// define http server
	server := &http.Server{
		Addr:    ":8080", // address to run on localHost or on google cloud server
		Handler: mux,
	}

	log.Println("Server listening on port 8080...")
	error := server.ListenAndServe()

	if errors.Is(error, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if error != nil {
		fmt.Printf("error listening for server: %s\n", error)
	}

	<-ctx.Done() // waiting indefinitely for context to be cancelled (should never happen)

	log.Println("Server Ended") // this should never be printed

}
