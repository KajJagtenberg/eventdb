package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dgraph-io/badger/v3"

	"github.com/eventflowdb/eventflowdb/env"
	"github.com/eventflowdb/eventflowdb/store"
	"github.com/eventflowdb/eventflowdb/transport"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})

	godotenv.Load()
}

func server() {
	log.Println("Starting EventflowDB")

	data := env.GetEnv("DATA", "data")

	var db *badger.DB
	var err error

	db, err = badger.Open(badger.DefaultOptions(data).WithLogger(log))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	eventstore, err := store.NewBadgerEventStore(store.BadgerStoreOptions{
		DB:             db,
		EstimateCounts: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer eventstore.Close()

	// grpcServer := transport.RunGRPCServer(eventstore, logger)
	// httpServer := transport.RunHTTPServer(eventstore, logger)
	restServer := transport.RunRestServer(eventstore, log)
	promServer := transport.RunPromServer(log)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	log.Println("EventflowDB is shutting down...")

	db.Close()
	// grpcServer.GracefulStop()
	// httpServer.Shutdown()
	restServer.Shutdown()
	promServer.Shutdown()
}

func main() {
	server()
}
