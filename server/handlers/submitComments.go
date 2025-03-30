package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type Comments struct {
	RecordingId int    `json:"recordingId"`
	Comments    string `json:"comments"`
}

func (env *Env) SubmitComments(w http.ResponseWriter, r *http.Request) {
	var c Comments
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Begin a database transaction
	tx, err := env.DB.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE recordings SET physician_comments = $2, status = 'viewed' WHERE recording_id = $1", c.RecordingId, c.Comments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("submitComments Response: %v\n", c)
}
