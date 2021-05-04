package main

import (
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/env"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/redcon"
	"go.etcd.io/bbolt"
)

var (
	data     = env.GetEnv("DATA", "data")
	port     = env.GetEnv("PORT", "6543")
	password = env.GetEnv("PASSWORD", "")

	log = logrus.New()
)

func check(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {
	log.SetFormatter(&logrus.JSONFormatter{})

	if len(password) == 0 {
		log.Warn("No password set")
	}

	log.Println("Initializing store")

	db, err := bbolt.Open(path.Join(data, "state.dat"), 0666, bbolt.DefaultOptions)
	check(err, "Failed to open database")
	defer db.Close()

	store, err := store.NewBoltStore(db, log)
	check(err, "Failed to create store")
	defer store.Close()

	log.Println("Initializing RESP server")

	go func() {
		log.Printf("RESP API listening on %s", port)

		commandHandler := api.Combine(
			api.AssertSession(),
			api.Authentication(password),
			api.CommandHandler(store),
		)

		acceptHandler := api.AcceptHandler()

		errorHandler := api.ErrorHandler()

		server := redcon.NewServer(":"+port, commandHandler, acceptHandler, errorHandler)

		check(server.ListenAndServe(), "Failed to run RESP API")
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	<-c

	log.Println("EventflowDB is shutting down...")
}
