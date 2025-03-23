package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// Define structs needed to parse through layers of Clerk Request Body
type ClerkReq struct {
	Physician Physician `json:"data"`
}

type Physician struct {
	ID             string         `json:"id"`
	EmailAddresses []EmailAddress `json:"email_addresses"`
}

type EmailAddress struct {
	EmailAddress string `json:"email_address"`
}

func (env *Env) CreatePhysician(w http.ResponseWriter, r *http.Request) {
	var d ClerkReq

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := d.Physician.ID
	email := d.Physician.EmailAddresses[0].EmailAddress

	if id == "" || email == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		log.Printf("createPhysician: %s %s\n", id, email)
		return
	}

	// Begin a database transaction
	tx, err := env.DB.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT INTO physician (physician_id, email) VALUES ($1, $2) ON CONFLICT (physician_id) DO UPDATE SET email = $2", id, email)
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
	log.Printf("createPhysician Response: %s %s\n", id, email)
}
