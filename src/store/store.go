package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"

	// ok
	_ "github.com/mattn/go-sqlite3"
)

// Roster provides sql access to Class db
type Roster struct {
	*sql.DB
}

// New returns rosters form provided db
func New(dbName string) (*Roster, error) {
	dir := filepath.Join(UserHomeDir(), "data")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0666)
	}
	db, err := sql.Open("sqlite3", filepath.Join(dir, dbName))
	if err != nil {
		return nil, err
	}
	return &Roster{db}, nil
}

// UserHomeDir returns the windows/unix homedir: Store rosters.db in home/data/rosters.db
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
