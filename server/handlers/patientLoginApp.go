package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type patientLogin struct {
	PatientId uint64 `db:"patient_id" json:"patientId"`
}

func (env *Env) PatientLoginApp(w http.ResponseWriter, r *http.Request) {

	// ensure receiving GET request
	if r.Method != "GET" {
		http.Error(w, "Invalid http request type", http.StatusBadRequest)
		log.Println("invalid http request type - should be GET request - instead is", r.Method)
	}

	patient := patientLogin{}

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

	attemptedUsername := q.Get("username") // pull username from url parameters
	attemptedPassword := q.Get("password") // pull password from url parameters

	// verify URL contained required inputs
	if attemptedUsername == "" || attemptedPassword == "" {
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

	// select patient id based on correct username/password from database
	err = tx.Get(&patient, `SELECT a.patient_id FROM app_login a WHERE (username, password) = ($1, $2)`, attemptedUsername, attemptedPassword)
	if err != nil {
		log.Printf("Error fetching data from database: %v\n", err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(patient); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("LoadRecordingInfoApp: Error encoding JSON: %v\n", err)
		return
	}
	log.Printf("patient response: %v\n", patient)

}
