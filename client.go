package main

import (
	"context"

	"github.com/kajjagtenberg/eventflowdb/transport"
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

	store := transport.NewEventStoreServiceClient(conn)
	res, err := store.Checksum(context.Background(), &transport.ChecksumRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.Checksum)
}
