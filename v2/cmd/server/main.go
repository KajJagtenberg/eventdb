package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Hello, world!")

	addr := ":6543"

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Unable to listen on %s", addr)
	}
	defer lis.Close()

	srv := grpc.NewServer()

	log.Println("Initializing gRPC services")

	log.Printf("Starting gRPC server on %s", addr)

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Unable to start gRPC server: %v", err)
	}
}
