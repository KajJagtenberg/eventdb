package api

import (
	"github.com/hashicorp/raft"
	"github.com/kajjagtenberg/eventflowdb/persistence"
	"google.golang.org/grpc"
)

func NewGRPCServer(raft *raft.Raft, persistence *persistence.Persistence) *grpc.Server {
	server := grpc.NewServer()

	RegisterStreamServiceServer(server, NewStreamService(raft, persistence))
	RegisterShellServiceServer(server, NewShellService(raft, persistence))
	RegisterRaftServiceServer(server, NewRaftService(raft))

	return server
}
