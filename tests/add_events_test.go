package tests

import (
	"eventdb/store"
	"eventdb/util"
	"log"
	"testing"

	"go.etcd.io/bbolt"
)

func TestAddMockEvents(t *testing.T) {
	db, err := bbolt.Open("../data.bolt", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	eventstore, err := store.NewStore(db)
	if err != nil {
		t.Fatal(err)
	}

	util.AddMockEvents(eventstore, 1000)
}
