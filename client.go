package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	conn, err := grpc.Dial("127.0.0.1:6543", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := api.NewEventStoreServiceClient(conn)

	for {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		_, err := store.Add(ctx, &api.AddRequest{
			Stream:  uuid.New().String(),
			Version: 0,
			Events: []*api.AddRequest_EventData{
				{
					Type: "TestEvent",
					Data: []byte(""),
				},
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(10 * time.Millisecond)
	}
}
