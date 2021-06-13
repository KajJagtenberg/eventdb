package main

import (
	"context"
	"encoding/json"
	"os"

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

	res, err := store.GetAll(context.Background(), &api.GetAllRequest{})
	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	enc.Encode(res.Events)
}
