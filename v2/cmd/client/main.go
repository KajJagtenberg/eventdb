package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/api"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:6543", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	streamService := api.NewStreamServiceClient(conn)

	stream := uuid.New()

	request := &api.AddEventsRequest{
		Stream:  stream[:],
		Version: 0,
		Events: []*api.AddEventsRequest_EventData{
			{
				Type: "TestEvent",
			},
		},
	}

	response, err := streamService.AddEvents(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to execute request: %v", err)
	}

	log.Println(response)
}
