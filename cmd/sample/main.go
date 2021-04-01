package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/util"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:6543", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	service := api.NewStreamServiceClient(conn)

	for event := range util.GenerateEvents(100000) {

		stream := uuid.New()

		request := &api.AddEventsRequest{
			Stream:  stream[:],
			Version: 0,
			Events:  []*api.AddEventsRequest_EventData{&event},
		}

		log.Printf("Adding event to %v", stream)

		_, err := service.AddEvents(context.Background(), request)
		if err != nil {
			log.Fatalf("Failed to add events: %v", err)
		}
	}

	log.Println("Done")
}
