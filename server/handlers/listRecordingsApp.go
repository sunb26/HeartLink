package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type recordingList struct {
	RecordingId       uint64 `db:"recording_id" json:"id"` // this json tag has to be different because of sql limitation
	RecordingDateTime string `db:"recording_datetime" json:"recordingDateTime"`
}

func (env *Env) ListRecordingsApp(w http.ResponseWriter, r *http.Request) {

	fmt.Print("ListRecordingsApp Endpoint - Start\n")

	// ensure receiving GET request
	if r.Method != "GET" {
		http.Error(w, "Invalid http request type", http.StatusBadRequest)
		log.Println("invalid http request type - should be GET request - instead is", r.Method)
	}

	newRecording := []recordingList{}

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

	patientId := q.Get("patientid") // pull patient id from url parameters

	// verify URL contained required inputs
	if patientId == "" {
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

	// select recording id and recording date/time from database
	err = tx.Select(&newRecording,
		`SELECT
		r.recording_id, 
		r.recording_datetime
	FROM
		recordings r
	WHERE patient_id = $1
	ORDER BY r.recording_id DESC`, patientId)
	if err != nil {
		log.Printf("Error fetching data from database: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var parsedTime time.Time
	var parsedTimeLocal time.Time

	location, err := time.LoadLocation("America/New_York") // set time zone to EST
	if err != nil {
		log.Printf("Error loading location: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update recordingDateTime format
	for j := range newRecording {
		parsedTime, err = time.Parse(time.RFC3339, newRecording[j].RecordingDateTime) // convert to time.Time general format
		if err != nil {
			log.Printf("Error parsing date/time: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		parsedTimeLocal = parsedTime.In(location)                                         // convert to EST time zone
		newRecording[j].RecordingDateTime = parsedTimeLocal.Format("2006-01-02 15:04:05") // format to YYYY-MM-DD HH:MM:SS
	}

	// TESTING
	for i := 0; i < len(newRecording); i++ {
		fmt.Printf("recording id: %d\n", newRecording[i].RecordingId)
		fmt.Printf("date/time: %s\n", newRecording[i].RecordingDateTime)
	}

	// create JSON response
	data := make(map[string]interface{})
	data["widgets"] = newRecording

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("LoadRecordingInfoApp: Error encoding JSON: %v\n", err)
		return
	}
	log.Printf("LoadRecordingInfoApp Response: %v\n", newRecording)

}
