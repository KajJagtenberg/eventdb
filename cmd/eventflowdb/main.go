package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/joho/godotenv"
	"github.com/kajjagtenberg/eventflowdb/env"
	"github.com/kajjagtenberg/eventflowdb/fsm"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/kajjagtenberg/eventflowdb/transport"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	log = logrus.New()
)

const (
	// The maxPool controls how many connections we will pool.
	maxPool = 3

	// The timeout is used to apply I/O deadlines. For InstallSnapshot, we multiply
	// the timeout by (SnapshotSize / TimeoutScale).
	// https://github.com/hashicorp/raft/blob/v1.1.2/net_transport.go#L177-L181
	tcpTimeout = 10 * time.Second

	// The `retain` parameter controls how many
	// snapshots are retained. Must be at least 1.
	raftSnapShotRetain = 2

	// raftLogCacheSize is the maximum number of logs to cache in-memory.
	// This is used to reduce disk I/O for the recently committed entries.
	raftLogCacheSize = 512
)

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})

	godotenv.Load()
}

func server() {
	data := env.GetEnv("DATA", "data")
	port := env.GetEnv("PORT", "6543")
	tlsEnabled := env.GetEnv("TLS_ENABLED", "false") == "true"
	certFile := env.GetEnv("TLS_CERT_FILE", "certs/cert.pem")
	keyFile := env.GetEnv("TLS_KEY_FILE", "certs/key.pem")

	log.Println("initializing store")

	db, err := badger.Open(badger.DefaultOptions(path.Join(data)).WithLogger(log))
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

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	if tlsEnabled {
		crt, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatal(err)
		}

		config := &tls.Config{
			Certificates: []tls.Certificate{crt},
		}

		lis = tls.NewListener(lis, config)
		if err != nil {
			log.Fatal(err)
		}
	}

	grpcServer := grpc.NewServer()

	transport.RegisterEventStoreServiceServer(grpcServer, transport.NewEventStoreService(eventstore))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	<-c

	log.Println("eventflowDB is shutting down...")
}

func testRaft() {
	db, err := badger.Open(badger.DefaultOptions("data/fsm"))
	if err != nil {
		log.Fatal(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	raftConf := raft.DefaultConfig()
	raftConf.SnapshotThreshold = 1024
	raftConf.LocalID = raft.ServerID(hostname)

	fsmStore := fsm.NewBadgerFSM(db)

	store, err := raftboltdb.NewBoltStore("data/store")
	if err != nil {
		log.Fatal(err)
	}

	cacheStore, err := raft.NewLogCache(raftLogCacheSize, store)
	if err != nil {
		log.Fatal(err)
	}

	snapshotStore, err := raft.NewFileSnapshotStore("data/snapshots", raftSnapShotRetain, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	raftBindAddr := "127.0.0.1:16543"

	tcpAddr, err := net.ResolveTCPAddr("tcp", raftBindAddr)
	if err != nil {
		log.Fatal(err)
	}

	transport, err := raft.NewTCPTransport(raftBindAddr, tcpAddr, maxPool, tcpTimeout, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	raftSrv, err := raft.NewRaft(raftConf, fsmStore, cacheStore, store, snapshotStore, transport)
	if err != nil {
		log.Fatal(err)
	}
	defer raftSrv.Shutdown()

	app := fiber.New()

	app.Post("/raft/join", func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	})

	app.Post("/raft/remove", func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	})

	app.Get("/raft/stats", func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	})

	app.Post("/raft/join", func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	})

	app.Post("/store/set", func(c *fiber.Ctx) error {
		var body struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}

		if err := c.BodyParser(body); err != nil {
			return err
		}

		if len(body.Key) == 0 {
			return errors.New("key cannot be empty")
		}

		if len(body.Value) == 0 {
			return errors.New("value cannot be empty")
		}

		cmd, err := json.Marshal(fsm.CommandPayload{
			Operation: "SET",
			Key:       body.Key,
			Value:     body.Value,
		})
		if err != nil {
			return err
		}

		applyFuture := raftSrv.Apply(cmd, time.Millisecond*500)
		if err := applyFuture.Error(); err != nil {
			return err
		}

		_, ok := applyFuture.Response().(*fsm.ApplyResponse)
		if !ok {
			return errors.New("invalid return value")
		}

		return c.SendString("key set")
	})

	app.Get("/store/get", func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	})

	app.Post("/store/delete", func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	})

	app.Listen(":3000")
}

func main() {
	// server()
	testRaft()
}
