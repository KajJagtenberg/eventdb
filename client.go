package main

import (
	"context"
	"math/rand"
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

		data := make([]byte, 150)
		metadata := make([]byte, 50)
		rand.Read(data)
		rand.Read(metadata)

		_, err := store.Add(ctx, &api.AddRequest{
			Stream:  uuid.New().String(),
			Version: 0,
			Events: []*api.AddRequest_EventData{
				{
					Type:     "TestEvent",
					Data:     data,
					Metadata: metadata,
				},
			},
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
