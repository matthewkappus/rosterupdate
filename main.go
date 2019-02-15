package main

import (
	"log"
	"time"

	"github.com/matthewkappus/rosterUpdate/src/store"
)

func main() {
	r, err := store.New("rosters.db")
	if err != nil {
		log.Fatal(err)
	}

	if err = r.DownloadRosters(time.Minute * 2); err != nil {
		log.Fatal(err)
	}
	println("completed")
}
