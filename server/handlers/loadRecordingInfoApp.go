package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type recordingInfo struct {
	Status            string `db:"status" json:"status"`
	PhysicianComments string `db:"physician_comments" json:"physicianComments"`
	RecordingId       uint64 `db:"recording_id" json:"recordingId"`
	DownloadUrl       string `db:"download_url" json:"downloadUrl"`
}

func (env *Env) LoadRecordingInfoApp(w http.ResponseWriter, r *http.Request) {

	// ensure receiving GET request
	if r.Method != "GET" {
		http.Error(w, "Invalid http request type", http.StatusBadRequest)
		log.Println("invalid http request type - should be GET request - instead is", r.Method)
	}

	newRecording := recordingInfo{}

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

	// verify URL contained required inputs
	if recordingId == "" {
		http.Error(w, "missing required URL inputs", http.StatusBadRequest)
		return
	}

	// start transaction
	tx, err := env.DB.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// select view status and physician comments from database
	err = tx.Get(&newRecording,
		`SELECT
		r.status,
		r.physician_comments,
		r.recording_id, 
		r.download_url
	FROM
		recordings r
	WHERE recording_id = $1`, recordingId)
	if err != nil {
		log.Printf("Error fetching data from database: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(newRecording); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("LoadRecordingInfoApp: Error encoding JSON: %v\n", err)
		return
	}
	log.Printf("LoadRecordingInfoApp Response: %v\n", newRecording)

}
