package api

import (
	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
)

func NewGRPCServer(raft *raft.Raft) *grpc.Server {
	streamService := NewStreamService(raft)

	server := grpc.NewServer()

	RegisterStreamServiceServer(server, streamService)

	return server
}
