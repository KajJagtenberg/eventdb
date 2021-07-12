package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/eventflowdb/eventflowdb/constants"
	"github.com/eventflowdb/eventflowdb/env"
	"github.com/eventflowdb/eventflowdb/store"
	"github.com/eventflowdb/eventflowdb/transport"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	db, err := gorm.Open(postgres.Open(env.GetEnv("DATABASE_URL", "postgresql://postgres:password@127.0.0.1:5432/eventflowdb")))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to Postgres")

	eventstore, err := store.NewSQLEventStore(db)
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
