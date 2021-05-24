package store_test

import (
	"crypto/rand"
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
)

func BenchmarkAddMemory(t *testing.B) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	eventstore, err := store.NewBadgerEventStore(store.BadgerStoreOptions{
		DB:             db,
		EstimateCounts: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < t.N; i++ {
		for j := 0; j < 10000; j++ {
			data := make([]byte, 150)
			metadata := make([]byte, 50)
			rand.Read(data)
			rand.Read(metadata)

			var events []store.EventData
			events = append(events, store.EventData{
				Type:     "AccountOpened",
				Data:     data,
				Metadata: metadata,
			})

			_, err := eventstore.Add(uuid.New(), 0, events)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
