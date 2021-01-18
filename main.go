package main

import (
	"log"

	"github.com/dgraph-io/badger/v3"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("data").WithLogger(nil))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := NewStore(db)

	// TODO: Add an networked API
}
