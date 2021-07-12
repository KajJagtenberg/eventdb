package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/eventflowdb/eventflowdb/constants"
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
	log.Printf("Starting EventflowDB v%s", constants.Version)

	eventstore, err := store.NewSQLStore()
	if err != nil {
		log.Fatal(err)
	}
	defer eventstore.Close()

	// grpcServer := transport.RunGRPCServer(eventstore, logger)
	restServer := transport.RunRestServer(eventstore, log)
	promServer := transport.RunPromServer(log)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	log.Println("EventflowDB is shutting down...")

	// grpcServer.GracefulStop()
	restServer.Shutdown()
	promServer.Shutdown()
}

func main() {
	server()
}
