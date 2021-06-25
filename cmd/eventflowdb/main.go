package main

import (
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/dgraph-io/badger/v3"

	"github.com/eventflowdb/eventflowdb/env"
	"github.com/eventflowdb/eventflowdb/store"
	"github.com/eventflowdb/eventflowdb/transport"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()
)

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})

	godotenv.Load()
}

func server() {
	data := env.GetEnv("DATA", "data")

	var db *badger.DB
	var err error

	db, err = badger.Open(badger.DefaultOptions(path.Join(data, "fsm")).WithLogger(nil).WithInMemory(false))
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	eventstore, err := store.NewBadgerEventStore(store.BadgerStoreOptions{
		DB:             db,
		EstimateCounts: true,
	})
	if err != nil {
		logger.Fatal(err)
	}
	defer eventstore.Close()

	grpcServer := transport.RunGRPCServer(eventstore, logger)
	httpServer := transport.RunHTTPServer(eventstore, logger)
	promServer := transport.RunPromServer(logger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	logger.Println("eventflowDB is shutting down...")

	db.Close()
	grpcServer.GracefulStop()
	httpServer.Shutdown()
	promServer.Shutdown()
}

func main() {
	server()
}
