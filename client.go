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
		stream := uuid.New()

		var events []*api.AddRequest_EventData

		for i := 0; i < 10; i++ {
			events = append(events, &api.AddRequest_EventData{
				Type:     "TestEvent",
				Data:     []byte("data"),
				Metadata: []byte("metadata"),
			})
		}

		res, err := store.Add(context.Background(), &api.AddRequest{
			Stream:  stream.String(),
			Version: 0,
			Events:  events,
		})
		if err != nil {
			log.Fatal(err)
		}

		log.Println(len(res.Events))

		time.Sleep(time.Millisecond * 10)
	}

}
