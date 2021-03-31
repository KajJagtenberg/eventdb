package main

import (
	"context"
	"log"
	"sync"

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

	var wg sync.WaitGroup

	for event := range util.GenerateEvents(1000) {
		wg.Add(1)

		stream := uuid.New()

		request := &api.AddEventsRequest{
			Stream:  stream[:],
			Version: 0,
			Events:  []*api.AddEventsRequest_EventData{&event},
		}

		log.Printf("Adding event to %v", stream)

		go func() {
			_, err := service.AddEvents(context.Background(), request)
			if err != nil {
				log.Fatalf("Failed to add events: %v", err)
			}
			wg.Done()
		}()
	}

	log.Println("Waiting")

	wg.Wait()

	log.Println("Done")
}
