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
	db, err := badger.Open(badger.DefaultOptions("data").WithLogger(nil))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	eventstore, err := store.NewBadgerEventStore(store.BadgerStoreOptions{
		DB: db,
	})
	if err != nil {
		log.Fatal(err)
	}

	events := []*api.EventData{}

	for i := 0; i < 25; i++ {
		events = append(events, &api.EventData{
			Id:   uuid.NewString(),
			Type: "TestEvent",
			Data: []byte("data"),
		})
	}

	res, err := eventstore.AppendToStream(&api.AppendToStreamRequest{
		Stream:  uuid.NewString(),
		Version: 0,
		Events:  events,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.Events)
}
