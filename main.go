package main

import (
	"log"
	"os"
	"time"

	"github.com/matthewkappus/rosterUpdate/src/store"
)

func main() {

	u := os.Getenv("SYNERGY_USER")
	p := os.Getenv("SYNERGY_PASSWORD")

	f, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	logger := log.New(f, "UpdateLog: ", log.Ldate|log.Lshortfile)
	rosterDB, err := store.New("/rosters.db")

	if err = rosterDB.DownloadRosters(time.Minute*2, u, p); err != nil {
		logger.Fatal(err)
	}

	logger.Print("Updated rosters")
}
