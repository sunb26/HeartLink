// main.go file

package main

import (
	"heartlinkServer/handlers"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s\n", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
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

	log.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", logging(mux)))
}
