package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dgraph-io/badger/v3"

	"github.com/eventflowdb/eventflowdb/env"
	"github.com/eventflowdb/eventflowdb/store"
	"github.com/eventflowdb/eventflowdb/transport"
	"github.com/gocql/gocql"
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
	restServer := transport.RunRestServer(eventstore, log)
	promServer := transport.RunPromServer(log)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	log.Println("EventflowDB is shutting down...")

	db.Close()
	// grpcServer.GracefulStop()
	restServer.Shutdown()
	promServer.Shutdown()
}

func main() {
	// server()

	cluster := gocql.NewCluster(strings.Split(env.GetEnv("CASSANDRA_NODES", "127.0.0.1"), ",")...) // TODO: Add to README
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: env.GetEnv("CASSANDRA_USERNAME", "cassandra"), // TODO: Add to README
		Password: env.GetEnv("CASSANDRA_PASSWORD", "cassandra"), // TODO: Add to README
	}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	log.Println("Connected to Cassandra cluster")

	session.SetConsistency(gocql.Quorum)

	queries := []*gocql.Query{
		session.Query("CREATE KEYSPACE IF NOT EXISTS ks1 WITH REPLICATION = {'class':'SimpleStrategy','replication_factor':3}"),
	}

	for _, qry := range queries {
		if err := qry.Exec(); err != nil {
			log.Fatal(err)
		}
	}
}
