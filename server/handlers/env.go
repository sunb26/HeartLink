package handlers

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Define dependencies for the handlers

type Env struct {
	DB *sqlx.DB
}
