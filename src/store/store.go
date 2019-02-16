package store

import (
	"database/sql"

	// ok
	_ "github.com/mattn/go-sqlite3"
)

// Roster provides sql access to Class db
type Roster struct {
	*sql.DB
}

// New returns rosters form provided db
func New(dbPath string) (*Roster, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &Roster{db}, nil
}
