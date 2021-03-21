package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
	"google.golang.org/grpc"
)

func main() {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	conn, err := grpc.Dial(":6543", grpc.WithInsecure(), grpc.WithTimeout(time.Second*10))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()

	c := store.NewStreamsClient(conn)

	stream := uuid.New()

	result, err := c.Add(ctx, &store.AddRequest{
		Stream:  stream[:],
		Version: 0,
		Events: []*store.AddRequest_Event{
			{
				Type:     "TestEvent",
				Data:     []byte("abcdefghijklmnopqrstuvw"),
				Metadata: []byte("metadata"),
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to add: %v", err)
	}

	for {
		m, err := result.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read message: %v", err)
		}

		log.Println(m.String())
	}
}
