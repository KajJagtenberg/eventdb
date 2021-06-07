package main

import (
	"context"

	"github.com/kajjagtenberg/eventflowdb/transport"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	client, err := grpc.Dial("127.0.0.1:6543", grpc.WithInsecure(), grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(metadata.AppendToOutgoingContext(ctx, "token", "this is my token"), method, req, reply, cc)
	}))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	store := transport.NewEventStoreServiceClient(client)

	ctx := context.Background()

	result, err := store.Checksum(ctx, &transport.ChecksumRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result.Checksum)
}
