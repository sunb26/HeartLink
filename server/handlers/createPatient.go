package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type Patient struct {
	PhysicianId string `json:"physicianId"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Dob         string `json:"dob"`
	Sex         string `json:"sex"`
	Height      int    `json:"height"`
	Weight      int    `json:"weight"`
}

func (env *Env) CreatePatient(w http.ResponseWriter, r *http.Request) {
	var p Patient
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

	_, err = tx.Exec("INSERT INTO patient (physician_id, firstname, lastname, email, height, weight, dob, sex, last_updated) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW()) ON CONFLICT (patient_id) DO UPDATE SET firstname = $2, lastname = $3, email = $4, height = $5, weight = $6, dob = $7, sex = $8, last_updated = NOW()", p.PhysicianId, p.FirstName, p.LastName, p.Email, p.Height, p.Weight, p.Dob, p.Sex)
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
	log.Printf("createPatient Response: %v\n", p)
}
