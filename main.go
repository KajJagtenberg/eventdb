package main

import (
	"log"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("data").WithLogger(nil))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := NewStore(db)

	cause := uuid.New().String()

	events := []AppendEvent{
		{
			Type: "product_added",
			Data: struct {
				Name  string `json:"name"`
				Price int    `json:"price"`
			}{
				"appel",
				400,
			},
			CausationID:   cause,
			CorrelationID: cause,
		},
	}

	if err := store.AppendToStream("mystream", 0, events); err != nil {
		log.Println(err)
	}

	loaded, err := store.LoadFromStream("mystream", 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(loaded[0])
}
