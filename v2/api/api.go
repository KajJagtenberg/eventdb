package api

import "google.golang.org/grpc"

func NewGRPCServer() *grpc.Server {
	streamService := NewStreamService()

	server := grpc.NewServer()

	RegisterStreamServiceServer(server, streamService)

	return server
}
