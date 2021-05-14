package main

import (
	"crypto/rand"
	"encoding/base32"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/kajjagtenberg/eventflowdb/commands"
	"github.com/kajjagtenberg/eventflowdb/env"
	"github.com/kajjagtenberg/eventflowdb/resp"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/redcon"
	"go.etcd.io/bbolt"
)

var (
	data       = env.GetEnv("DATA", "data")
	port       = env.GetEnv("PORT", "6543")
	password   = env.GetEnv("PASSWORD", "")
	noPassword = env.GetEnv("NO_PASSWORD", "false") == "true"

	log = logrus.New()
)

func check(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {
	log.SetFormatter(&logrus.JSONFormatter{})

	if !noPassword && len(password) == 0 {
		passwordData := make([]byte, 20)
		_, err := rand.Read(passwordData)
		check(err, "Failed to generate password")

		password = base32.StdEncoding.EncodeToString(passwordData)

		log.Printf("Generated a password since none was given: %s", password)
	}

	log.Println("Initializing store")

	db, err := bbolt.Open(path.Join(data, "state.dat"), 0666, bbolt.DefaultOptions)
	check(err, "Failed to open database")
	defer db.Close()

	eventstore, err := store.NewBoltStore(db, log)
	check(err, "Failed to create store")
	defer eventstore.Close()

	dispatcher := commands.NewCommandDispatcher()
	dispatcher.Register("uptime", "u", commands.UptimeHandler())
	dispatcher.Register("version", "v", commands.VersionHandler())
	dispatcher.Register("add", "a", commands.AddHandler(eventstore))
	dispatcher.Register("checksum", "ch", commands.ChecksumHandler(eventstore))
	dispatcher.Register("eventcount", "ec", commands.EventCountHandler(eventstore))
	dispatcher.Register("eventcountest", "ece", commands.EventCountEstimateHandler(eventstore))
	dispatcher.Register("get", "g", commands.GetHandler(eventstore))
	dispatcher.Register("getall", "ga", commands.GetHandler(eventstore))

	log.Println("Initializing RESP server")

	go func() {
		log.Printf("RESP API listening on %s", port)

		server := redcon.NewServer(":"+port, resp.CommandHandler(dispatcher), resp.AcceptHandler(), resp.ErrorHandler())

		check(server.ListenAndServe(), "Failed to run RESP API")
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	<-c

	log.Println("EventflowDB is shutting down...")
}
