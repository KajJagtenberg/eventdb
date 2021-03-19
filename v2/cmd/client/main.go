package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":6543", grpc.WithInsecure(), grpc.WithTimeout(time.Second*10))
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()

	c := store.NewEventStoreClient(conn)

	stream, err := uuid.New().MarshalBinary()
	if err != nil {
		log.Fatalf("Unable to marshal stream: %v", err)
	}

	req := &store.AddRequest{
		Stream:  stream,
		Version: 0,
		Events: []*store.AddRequest_Event{
			{
				Type: "ProductAdded",
				Data: []byte(`{"name":"Samsung Galaxy S8","version":80000}`),
			},
		},
	}

	res, err := c.Add(context.Background(), req)
	if err != nil {
		log.Fatalf("Unable to perform request: %v", err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(res)
}
