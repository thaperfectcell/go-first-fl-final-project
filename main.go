package main

import (
	"log"
	"os"

	"github.com/thaperfectcell/go-first-fl-final-project/pkg/db"
	"github.com/thaperfectcell/go-first-fl-final-project/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
