package handlers

import (
	"encoding/json"
	"net/http"
)

type Physician struct {
	Id    string
	Email string
}

func (env *Env) CreatePhysician(w http.ResponseWriter, r *http.Request) {
	var p Physician

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&p)
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

	_, err = tx.Exec("INSERT INTO physician (physician_id, email) VALUES ($1, $2)", p.Id, p.Email)
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
}
