package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/google/uuid"

	"github.com/eventflowdb/eventflowdb/constants"
	"github.com/eventflowdb/eventflowdb/env"
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
	log.Printf("Starting EventflowDB v%s", constants.Version)

	data := env.GetEnv("DATA", "data")

	var db *badger.DB
	var err error

	db, err = badger.Open(badger.DefaultOptions(data).WithLogger(log))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// eventstore, err := store.NewBadgerEventStore(store.BadgerStoreOptions{
	// 	DB:             db,
	// 	EstimateCounts: true,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer eventstore.Close()

	// grpcServer := transport.RunGRPCServer(eventstore, logger)
	// restServer := transport.RunRestServer(eventstore, log)
	promServer := transport.RunPromServer(log)

	app := fiber.New()
	app.Use(compress.New())
	app.Get("/stream/:stream", func(c *fiber.Ctx) error {
		stream, err := uuid.Parse(c.Params("stream"))
		if err != nil {
			return err
		}

		offset, err := strconv.ParseInt(c.Query("offset", "0"), 10, 32)
		if err != nil {
			return err
		}

		txn := db.NewTransaction(false)
		defer txn.Discard()

		key := new(bytes.Buffer)
		key.WriteByte(0)
		key.WriteString(stream.String())

		binary.Write(key, binary.BigEndian, offset)

		item, err := txn.Get(key.Bytes())
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return c.Send(val)
		})
	})

	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	log.Println("EventflowDB is shutting down...")

	db.Close()
	// grpcServer.GracefulStop()
	// restServer.Shutdown()
	promServer.Shutdown()
}

func cassandra() {
	cluster := gocql.NewCluster(strings.Split(env.GetEnv("CASSANDRA_NODES", "127.0.0.1"), ",")...) // TODO: Add to README
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: env.GetEnv("CASSANDRA_USERNAME", "cassandra"), // TODO: Add to README
		Password: env.GetEnv("CASSANDRA_PASSWORD", "cassandra"), // TODO: Add to README
	}
	cluster.Keyspace = "ks1"

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	log.Println("Connected to Cassandra cluster")

	session.SetConsistency(gocql.Quorum)

	queries := []*gocql.Query{
		session.Query("DROP TABLE IF EXISTS events;"),
		session.Query("CREATE TABLE events (id uuid, event blob, PRIMARY KEY(id)) WITH compaction = {'class':'SizeTieredCompactionStrategy'};"),
	}

	for _, qry := range queries {
		if err := qry.Exec(); err != nil {
			log.Fatal(err)
		}
	}

	start := time.Now()

	b := session.NewBatch(gocql.LoggedBatch)

	for i := 0; i < 1000; i++ {
		b.Query("INSERT INTO events (id, event) VALUES (?, ?);", uuid.NewString(), []byte("hello there!"))
	}

	if err := session.ExecuteBatch(b); err != nil {
		log.Fatal(err)
	}

	log.Println(time.Since(start))

	iter := session.Query("SELECT * FROM events;").Consistency(gocql.One).Iter()
	defer iter.Close()

	var id string
	var blob []byte

	for iter.Scan(&id, &blob) {
		// log.Println(id, string(blob))
	}
}

func main() {
	server()
	// cassandra()
}
