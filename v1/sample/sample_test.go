package sample

import (
	"eventflowdb/store"
	"log"
	"testing"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

func TestAddSampleEvents(t *testing.T) {
	db, err := bbolt.Open("events.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	eventstore, err := store.NewEventStore(db)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10000; i++ {
		go func() {
			if _, err := eventstore.AppendToStream(uuid.New(), 0, []store.EventData{
				{
					Type: "ProductAdded",
					Data: []byte(`{"name":"Appeltaart","price":400}`),
				},
			}); err != nil {
				log.Fatal(err)
			}
		}()
	}
}
