package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type PatientPage struct {
	FirstName  string      `json:"firstName"`
	LastName   string      `json:"lastName"`
	Email      string      `json:"email"`
	Age        int         `json:"age"`
	Sex        string      `json:"sex"`
	Weight     int         `json:"weight"`
	Height     int         `json:"height"`
	Recordings []Recording `json:"recordings"`
}

type Recording struct {
	RecordingId       int    `db:"recording_id" json:"recordingId"`
	RecordingDateTime string `db:"recording_datetime" json:"recordingDateTime"`
	DownloadUrl       string `db:"download_url" json:"downloadUrl"`
	Comments          string `db:"comments" json:"comments"`
	HeartRate         int    `db:"heart_rate" json:"heartRate"`
}

func (env *Env) GetPatient(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("listPatients: %v\n", err)
		return
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("listPatients: %v\n", err)
		return
	}

	patientId := q.Get("patientId")

	if patientId == "" {
		http.Error(w, "getPatient: missing required fields", http.StatusBadRequest)
		return
	}

	recordings := []Recording{}
	p := PatientPage{}

	// Begin a database transaction
	tx, err := env.DB.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Get(&p, "SELECT firstname, lastname, email, DATE_PART('year', AGE(dob)) AS age, sex, weight, height FROM patient WHERE patient_id = $1", patientId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("getPatient: %v\n", err)
		return
	}

	err = tx.Select(&recordings, "SELECT recording_id, recording_datetime, download_url, physician_comments AS comments, heart_rate FROM recordings WHERE patient_id = $1 AND status <> 'notSubmitted' ORDER BY (recording_datetime) DESC", patientId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("getPatient: %v\n", err)
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
	for j := range recordings {
		parsedTime, err = time.Parse(time.RFC3339, recordings[j].RecordingDateTime) // convert to time.Time general format
		if err != nil {
			log.Printf("Error parsing date/time: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		parsedTimeLocal = parsedTime.In(location)                                       // convert to EST time zone
		recordings[j].RecordingDateTime = parsedTimeLocal.Format("2006-01-02 15:04:05") // format to YYYY-MM-DD HH:MM:SS
	}

	p.Recordings = recordings

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("listPatients: Error encoding JSON: %v\n", err)
		return
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("getPatient Response: %v\n", p)
}
