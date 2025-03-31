package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type patient struct {
	PatientId   int    `db:"patient_id" json:"patientId"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	LastUpdated string `db:"last_updated" json:"lastUpdated"`
	Verified    bool   `json:"verified"`
	Viewed      bool   `json:"viewed"`
}

func (env *Env) ListPatients(w http.ResponseWriter, r *http.Request) {
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

	physicianId := q.Get("physicianid")

	tx, err := env.DB.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("listPatients: %v\n", err)
		return
	}
	defer tx.Rollback()

	patients := []patient{}

	err = tx.Select(&patients,
		`SELECT 
    p.patient_id,
    p.firstname,
    p.lastname,
    p.email,
    p.last_updated,
    CASE WHEN al.patient_id IS NOT NULL THEN true ELSE false END AS verified,
    CASE WHEN COUNT(r.recording_id) FILTER (WHERE r.status != 'viewed') = 0 THEN true ELSE false END AS viewed
	FROM 
			patient p
	LEFT JOIN 
			app_login al ON p.patient_id = al.patient_id
	LEFT JOIN 
			recordings r ON p.patient_id = r.patient_id
	WHERE physician_id = $1
	GROUP BY p.patient_id, al.patient_id
	ORDER BY viewed`, physicianId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("listPatients: %v\n", err)
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
	for j := range patients {
		parsedTime, err = time.Parse(time.RFC3339, patients[j].LastUpdated) // convert to time.Time general format
		if err != nil {
			log.Printf("Error parsing date/time: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		parsedTimeLocal = parsedTime.In(location)                               // convert to EST time zone
		patients[j].LastUpdated = parsedTimeLocal.Format("2006-01-02 15:04:05") // format to YYYY-MM-DD HH:MM:SS
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["patients"] = patients

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("listPatients: Error encoding JSON: %v\n", err)
		return
	}
	log.Printf("listPatients Response: %v\n", patients)
}
