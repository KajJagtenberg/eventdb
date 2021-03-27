package api

import (
	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
)

func NewGRPCServer(raft *raft.Raft) *grpc.Server {

	server := grpc.NewServer()

	RegisterStreamServiceServer(server, NewStreamService(raft))
	RegisterShellServiceServer(server, NewShellService())

	return server
}
