package main

import (
	"context"

	service "github.com/kajjagtenberg/eventflowdb/grpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	client, err := grpc.Dial("127.0.0.1:6543", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	chat := service.NewChatClient(client)

	msg, err := chat.SendMessage(context.Background(), &service.Message{
		Body: "Hi there",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(msg.Body)
}
