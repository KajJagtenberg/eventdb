package main

import (
	"os"
	"os/signal"
	"path"
	"syscall"

	"log"

	"github.com/joho/godotenv"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/env"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/tidwall/redcon"
	"go.etcd.io/bbolt"

	_ "embed"
)

var (
	data     = env.GetEnv("DATA", "data/state.dat")
	respAddr = env.GetEnv("RESP_ADDR", ":6543")
	password = env.GetEnv("PASSWORD", "")
)

func main() {
	godotenv.Load()

	if len(password) == 0 {
		log.Println("WARNING: No password set")
	}

	log.Println("Initializing store")

	db, err := bbolt.Open(path.Join(data, "state.dat"), 0666, bbolt.DefaultOptions)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	store, err := store.NewBoltStore(db)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	log.Println("Initializing RESP server")

	go func() {
		log.Printf("RESP API listening on %s", respAddr)

		commandHandler := api.Combine(
			api.AssertSession(),
			api.Authentication(password),
			api.CommandHandler(store),
		)

		acceptHandler := api.AcceptHandler()

		errorHandler := api.ErrorHandler()

		if err := redcon.ListenAndServe(respAddr, commandHandler, acceptHandler, errorHandler); err != nil {
			log.Fatalf("Failed to run RESP API: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	<-c

	log.Println("EventflowDB is shutting down...")
}
