package main

import (
	"context"
	"log"

	"github.com/kajjagtenberg/eventflowdb/api"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:6543", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	streamService := api.NewStreamServiceClient(conn)

	request := &api.AddEventsRequest{}

	response, err := streamService.AddEvents(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to execute request: %v", err)
	}

	log.Println(response)
}
