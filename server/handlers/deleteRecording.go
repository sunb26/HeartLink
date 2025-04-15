package handlers

import (
	"log"
	"net/http"
	"net/url"
)

func (env *Env) DeleteRecording(w http.ResponseWriter, r *http.Request) {

	// ensure receiving DELETE request
	if r.Method != "DELETE" {
		http.Error(w, "Invalid http request type", http.StatusBadRequest)
		log.Println("invalid http request type - should be DELETE request - instead is", r.Method)
	}

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

	// delete record based on recording id from database
	result, err := tx.Exec("DELETE FROM recordings WHERE recording_id = $1", recordingId)
	if err != nil {
		log.Printf("Error fetching data from database: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// verify that row was deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "No rows affected", http.StatusNotFound)
		return
	}

}
