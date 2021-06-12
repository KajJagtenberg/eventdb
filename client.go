package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/transport"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func consumer() {
	client, err := grpc.Dial("127.0.0.1:6543", grpc.WithInsecure(), grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(metadata.AppendToOutgoingContext(ctx, "token", "7UAFSAQFKIVALAADRNB7GZKXMDUZPLTJV6RH2XA="), method, req, reply, cc)
	}))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	store := transport.NewEventStoreServiceClient(client)

	var checkpoint string

	for {
		log.Println(checkpoint)

		res, err := store.GetAll(context.Background(), &transport.GetAllRequest{
			Offset: checkpoint,
			Limit:  10,
		})
		if err != nil {
			log.Fatal(err)
		}

		l := len(res.Events)

		if l > 0 {
			checkpoint = res.Events[len(res.Events)-1].Id
		}
	}
}

func producer() {
	client, err := grpc.Dial("127.0.0.1:6543", grpc.WithInsecure(), grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(metadata.AppendToOutgoingContext(ctx, "token", "7UAFSAQFKIVALAADRNB7GZKXMDUZPLTJV6RH2XA="), method, req, reply, cc)
	}))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	store := transport.NewEventStoreServiceClient(client)

	for {
		_, err := store.Add(context.Background(), &transport.AddRequest{
			Stream:  uuid.New().String(),
			Version: 0,
			Events: []*transport.AddRequest_EventData{
				{
					Type:     "TestEvent",
					Data:     make([]byte, 10),
					Metadata: make([]byte, 10),
				},
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Microsecond * 100)
	}
}

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	go consumer()

	producer()
}
