package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/chzyer/readline"
	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/api"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:6543", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	func() {
		service := api.NewStreamServiceClient(conn)

		stream := uuid.New()

		addResponse, err := service.AddEvents(context.Background(), &api.AddEventsRequest{
			Stream:  stream[:],
			Version: 0,
			Events: []*api.AddEventsRequest_EventData{
				{
					Type:     "TestEvent",
					Data:     []byte("data"),
					Metadata: []byte("metadata"),
					AddedAt:  time.Now().UnixNano(),
				},
			},
		})
		if err != nil {
			log.Fatalf("Failed to add events: %v", err)
		}
		log.Println(addResponse.Events)

		getResponse, err := service.GetEvents(context.Background(), &api.GetEventsRequest{
			Stream:  stream[:],
			Version: 0,
			Limit:   0,
		})
		if err != nil {
			log.Fatalf("Failed to get events: %v", err)
		}
		log.Println(getResponse.Events)
	}()

	service := api.NewShellServiceClient(conn)

	rl, err := readline.New("> ")
	if err != nil {
		log.Fatalf("Failed to create readline: %v", err)
	}

	stream, err := service.Execute(context.Background())
	if err != nil {
		log.Fatalf("Failed to execute request: %v", err)
	}

	for {
		line, err := rl.Readline()
		if err == io.EOF {
			break
		}
		if err == readline.ErrInterrupt {
			break
		}

		if len(line) == 0 {
			continue
		}

		if err := stream.Send(&api.ShellRequest{
			Body: line,
		}); err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}

		response, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive response: %v", err)
		}

		if len(response.Body) == 0 {
			continue
		}

		fmt.Println(response.Body)
	}
}
