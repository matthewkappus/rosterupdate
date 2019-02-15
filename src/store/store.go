package store

import (
	"database/sql"

	// ok
	_ "github.com/mattn/go-sqlite3"
)

// Rosters provides sql access to Class db
type Rosters struct {
	*sql.DB
}

// New returns rosters form provided db
func New(dbPath string) (*Rosters, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &Rosters{db}, nil
}
