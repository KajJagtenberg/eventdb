package main

import (
	"os"
	"os/signal"
	"syscall"

	"log"

	"github.com/kajjagtenberg/eventflowdb/env"
	"github.com/kajjagtenberg/eventflowdb/store"
	"go.etcd.io/bbolt"
)

var (
	// localID       = env.GetEnv("RAFT_LOCAL_ID", "main")
	// bindAddr      = env.GetEnv("RAFT_BIND_ADDR", "127.0.0.1:6542")
	// advrAddr      = env.GetEnv("RAFT_ADVR_ADDR", bindAddr)
	// bootstrap     = env.GetEnv("RAFT_BOOTSTRAP", "false") == "true"
	stateLocation = env.GetEnv("STATE_LOCATION", "data/state.dat")
	grpcAddr      = env.GetEnv("GRPC_ADDR", ":6543")
	// graphqlAddr   = env.GetEnv("GRAPHQL_ADDR", ":16543")
)

func main() {
	db, err := bbolt.Open(stateLocation, 0666, bbolt.DefaultOptions)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	store, err := store.NewBoltStore(db)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	log.Println("Size: ", store.Size())

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
