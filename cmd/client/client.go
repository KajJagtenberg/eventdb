package main

import (
	"log"

	"github.com/google/uuid"
	v1 "github.com/kajjagtenberg/eventflowdb/client"
	"github.com/kajjagtenberg/eventflowdb/store"
)

func main() {
	client := v1.NewClient(&v1.Config{
		Address: "127.0.0.1:6543",
	})

	stream := uuid.New()

	_, err := client.Add(stream, 0, []store.EventData{
		{
			Type: "TestEvent",
			Data: []byte("This is the data"),
		},
	})
	if err != nil {
		log.Fatalf("Failed to add events: %v", err)
	}

	events, err := client.GetAll(v1.SINCE_START, 10)
	if err != nil {
		log.Fatalf("Failed to get events: %v", err)
	}

	log.Println(events)
}
