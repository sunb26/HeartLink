package main

import (
	"heartlinkServer/handlers"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE, PUT")
		w.Header().Add("Access-Control-Allow-Headers", "content-type")
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s\n", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {

	/* COMMENT OUT WHEN COMMITTING - FOR RUNNING WITH LOCAL HOST ONLY */
	// _, currentFile, _, ok := runtime.Caller(0) // get current file path
	// if !ok {
	// 	log.Fatalf("Unable to get current file info")
	// }
	// rootDir := filepath.Dir(filepath.Dir(currentFile))   // get root directory of current file
	// err := godotenv.Load(filepath.Join(rootDir, ".env")) // load environment variables from .env file
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }
	/* END COMMENT OUT WHEN COMMITTING - FOR RUNNING WITH LOCAL HOST ONLY */

	// connect to the SQL database
	db, err := sqlx.Connect("postgres", os.Getenv("SQL_DSN"))
	if err != nil {
		log.Fatalf("failed to open to database: %v", err)
	}

	// Instantiate dependencies to be injected for the handlers
	env := &handlers.Env{DB: db}

	mux := http.NewServeMux() // create custom multiplexer to handle incoming requests

	// CORS middleware
	corsMiddleware := cors.Default().Handler(mux)

	// each HandleFunc is used to handle a specific endpoint
	mux.HandleFunc("/UploadFilterRecording", env.UploadFilterRecording)
	mux.HandleFunc("/CreatePhysician", env.CreatePhysician)
	mux.HandleFunc("/ListPatients", env.ListPatients)
	mux.HandleFunc("/SaveRunAlgorithm", env.SaveRunAlgorithm)
	mux.HandleFunc("/LoadRecordingInfoApp", env.LoadRecordingInfoApp)
	mux.HandleFunc("/CreatePatient", env.CreatePatient)
	mux.HandleFunc("/SubmitComments", env.SubmitComments)
	mux.HandleFunc("/GetPatient", env.GetPatient)
	mux.HandleFunc("/DeleteRecording", env.DeleteRecording)

	log.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", logging(corsMiddleware)))
}
