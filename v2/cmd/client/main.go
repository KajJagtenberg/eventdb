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

	c := store.NewEventStoreClient(conn)

	result, err := c.GetStreams(ctx, &store.GetStreamsRequest{})
	if err != nil {
		log.Fatalf("Failed to perform request: %v", err)
	}

	for {
		stream, err := result.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive message: %v", err)
		}

		log.Println(uuid.FromBytes(stream.Id))
	}
}
