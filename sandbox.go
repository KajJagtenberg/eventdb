package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/api"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:6543", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := api.NewStreamServiceClient(conn)

	for {
		wg := sync.WaitGroup{}

		for i := 0; i < 10000; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				stream := uuid.New()

				data := make([]byte, 120)
				metadata := make([]byte, 40)

				rand.Read(data)
				rand.Read(metadata)

				response, err := client.AddEvents(context.Background(), &api.AddEventsRequest{
					Stream:  stream[:],
					Version: 0,
					Events: []*api.AddEventsRequest_EventData{
						{
							Type:     "RandomEvent",
							Data:     data,
							Metadata: metadata,
						},
					},
				})
				if err != nil {
					log.Fatalf("Failed to add events: %v", err)
				}

				for _, event := range response.Events {
					fmt.Println(event.Id)
				}
			}()
		}

		wg.Wait()
	}
}
