package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/matthewkappus/rosterUpdate/src/store"
)

var (
	u = flag.String("u", "", "Synergy User Name: Must have admin rights")
	p = flag.String("p", "", "Synergy Password: Must have admin rights")
)

func main() {
	flag.Parse()
	if *u == "" || *p == "" {
		log.Fatal("Must provide -u and -p flags")
	}

	f, err := os.OpenFile("log/rosterUpdate.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// TODO: create log/ data/ files if not exist
	logger := log.New(f, "UpdateLog: ", log.Ldate|log.Lshortfile)
	rosterDB, err := store.New("data/rosters.db")

	logger.Printf("Getting roster for %s", *u)
	if err = rosterDB.DownloadRosters(time.Minute*2, *u, *p); err != nil {
		logger.Fatal(err)
	} else {
		logger.Print("Updated rosters")
	}

}
