package main

import (
	"log"
	"net"

	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/cluster"
	"github.com/kajjagtenberg/eventflowdb/env"
	"github.com/kajjagtenberg/eventflowdb/persistence"
	"go.etcd.io/bbolt"
)

var (
	localID       = env.GetEnv("RAFT_LOCAL_ID", "")
	bindAddr      = env.GetEnv("RAFT_BIND_ADDR", ":6542")
	advrAddr      = env.GetEnv("RAFT_ADVR_ADDR", bindAddr)
	bootstrap     = env.GetEnv("RAFT_BOOTSTRAP", "false") == "true"
	stateLocation = env.GetEnv("STATE_LOCATION", "data/state.dat")
)

func main() {
	db, err := bbolt.Open(stateLocation, 0666, bbolt.DefaultOptions)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	persistence, err := persistence.NewPersistence(db)
	if err != nil {
		log.Fatalf("Failed to create persistence: %v", err)
	}

	fsm, err := cluster.NewFSM(persistence)
	if err != nil {
		log.Fatalf("Failed to create FSM: %v", err)
	}

	raftServer, err := cluster.NewRaftServer(localID, bindAddr, advrAddr, fsm, true)
	if err != nil {
		log.Fatalf("Failed to create Raft: %v", err)
	}
	defer raftServer.Shutdown()

	lis, err := net.Listen("tcp", ":6543")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	grpcServer := api.NewGRPCServer(raftServer, persistence)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
