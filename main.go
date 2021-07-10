package main

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/eventflowdb/eventflowdb/api"
	"github.com/eventflowdb/eventflowdb/store"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})

	godotenv.Load()
}

func main() {
	db, err := badger.Open(badger.DefaultOptions("data").WithLogger(log))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	eventstore, err := store.NewBadgerEventStore(store.BadgerStoreOptions{
		DB: db,
	})

	eventstore.AppendToStream(&api.AppendToStreamRequest{
		Stream:  uuid.NewString(),
		Version: 0,
		Events: []*api.EventData{
			{
				Id:   uuid.NewString(),
				Type: "TestEvent",
				Data: []byte("data"),
			},
		},
	})

	txn := db.NewTransaction(false)
	defer txn.Discard()

	iter := txn.NewIterator(badger.DefaultIteratorOptions)
	defer iter.Close()

	for iter.Rewind(); iter.Valid(); iter.Next() {
		item := iter.Item()

		log.Println(item.Key())
	}
}
