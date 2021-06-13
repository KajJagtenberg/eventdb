package main

import (
	"context"

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

	res, err := store.ClusterStats(context.Background(), &api.ClusterStatsRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.Stats)
}
