package main

import (
	"log"
	"net"

	"github.com/kajjagtenberg/eventflowdb/store"
	"go.etcd.io/bbolt"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Hello, world!")

	log.Println("Initializing Storage service")

	db, err := bbolt.Open("events.db", 0666, nil)
	if err != nil {
		log.Fatalf("Failed to initialize Storage service: %v", err)
	}
	defer db.Close()

	storage, err := store.NewStorage(db)
	if err != nil {
		log.Fatalf("Failed to initialize Storage service: %v", err)
	}

	addr := ":6543"

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s", addr)
	}
	defer lis.Close()

	srv := grpc.NewServer()

	log.Println("Initializing gRPC services")

	store.RegisterEventStoreServer(srv, store.NewStoreService(storage))

	log.Printf("Starting gRPC server on %s", addr)

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
