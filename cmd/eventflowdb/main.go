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
	"github.com/kajjagtenberg/eventflowdb/web"
	"github.com/kajjagtenberg/go-commando"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/redcon"
	"go.etcd.io/bbolt"
)

var (
	data       = env.GetEnv("DATA", "data")
	port       = env.GetEnv("PORT", "6543")
	httpPort   = env.GetEnv("HTTP_PORT", "16543")
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
		check(err, "failed to generate password")

		password = base32.StdEncoding.EncodeToString(passwordData)

		log.Printf("generated a password since none was given: %s", password)
	}

	log.Println("initializing store")

	db, err := bbolt.Open(path.Join(data, "state.dat"), 0666, bbolt.DefaultOptions)
	check(err, "failed to open database")
	defer db.Close()

	eventstore, err := store.NewBoltStore(db, log)
	check(err, "failed to create store")
	defer eventstore.Close()

	dispatcher := commando.NewCommandDispatcher()
	dispatcher.Register(commands.CMD_ADD, commands.CMD_ADD_SHORT, commands.AddHandler(eventstore))
	dispatcher.Register(commands.CMD_CHECKSUM, commands.CMD_CHECKSUM_SHORT, commands.ChecksumHandler(eventstore))
	dispatcher.Register(commands.CMD_EVENT_COUNT, commands.CMD_EVENT_COUNT_SHORT, commands.EventCountHandler(eventstore))
	dispatcher.Register(commands.CMD_EVENT_COUNT_EST, commands.CMD_EVENT_COUNT_EST_SHORT, commands.EventCountEstimateHandler(eventstore))
	dispatcher.Register(commands.CMD_GET, commands.CMD_GET_SHORT, commands.GetHandler(eventstore))
	dispatcher.Register(commands.CMD_GET_ALL, commands.CMD_GET_ALL_SHORT, commands.GetHandler(eventstore))
	dispatcher.Register(commands.CMD_PING, commands.CMD_PING_SHORT, commands.PingHandler())
	dispatcher.Register(commands.CMD_SIZE, commands.CMD_SIZE_SHORT, commands.SizeHandler(eventstore))
	dispatcher.Register(commands.CMD_STREAM_COUNT, commands.CMD_STREAM_COUNT_SHORT, commands.StreamCountHandler(eventstore))
	dispatcher.Register(commands.CMD_STREAM_COUNT_EST, commands.CMD_STREAM_COUNT_EST_SHORT, commands.StreamCountEstimateHandler(eventstore))
	dispatcher.Register(commands.CMD_UPTIME, commands.CMD_UPTIME_SHORT, commands.UptimeHandler())
	dispatcher.Register(commands.CMD_VERSION, commands.CMD_VERSION_SHORT, commands.VersionHandler())
	dispatcher.Register(commands.CMD_LIST_STREAMS, commands.CMD_LIST_STREAMS_SHORT, commands.ListStreamsHandler(eventstore))

	log.Println("initializing RESP server")

	go func() {
		log.Printf("RESP API listening on %s", port)

		server := redcon.NewServer(":"+port, resp.CommandHandler(dispatcher, password), resp.AcceptHandler(), resp.ErrorHandler())

		check(server.ListenAndServe(), "failed to run RESP API")
	}()

	app, err := web.CreateWebServer(dispatcher)
	check(err, "failed to create web server")

	go func() {
		log.Printf("HTTP API listening on %s", httpPort)

		check(app.Listen(":"+httpPort), "failed to run HTTP API")
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	<-c

	log.Println("eventflowDB is shutting down...")
}
