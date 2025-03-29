package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type recordingUrl struct {
	DownloadUrl string `db:"download_url" json:"downloadUrl"`
}

func (env *Env) StreamRecordingApp(w http.ResponseWriter, r *http.Request) {

	fmt.Print("StreamRecordingApp Endpoint - Start\n")

	// ensure receiving GET request
	if r.Method != "GET" {
		log.Println("invalid http request type - should be GET request - instead is", r.Method)
	}

	newRecording := recordingUrl{}

	// parse query parameter from URL
	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error parsing into URL structure: %v\n", err)
		return
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error parsing URL : %v\n", err)
		return
	}

	recordingId := q.Get("recordingid")

	fmt.Printf("recordingId: %s\n", recordingId) // TESTING

	// verify URL contained required inputs
	if recordingId == "" {
		http.Error(w, "missing required URL inputs", http.StatusBadRequest)
		return
	}

	// setup database connection
	tx, err := env.DB.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// get view downloadUrl from database
	err = tx.Get(&newRecording,
		`SELECT
		r.download_url
	FROM
		recordings r
	WHERE recording_id = $1`, recordingId)
	if err != nil {
		log.Printf("Error fetching status or physician comments from database: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Printf("downloadUrl: %s\n", newRecording.DownloadUrl) // TESTING

	// create JSON response
	data := make(map[string]interface{})
	data["recording"] = newRecording // SEE HOW BEN WANTS THIS LABELLED

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("StreamRecordingApp: Error encoding JSON: %v\n", err)
		return
	}
	log.Printf("StreamRecordingApp Response: %v\n", newRecording)

}
