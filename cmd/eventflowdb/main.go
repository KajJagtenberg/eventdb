package main

import (
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/dgraph-io/badger/v3"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
	promPort := env.GetEnv("PROM_PORT", "17654")
	tlsEnabled := env.GetEnv("TLS_ENABLED", "false") == "true"
	certFile := env.GetEnv("TLS_CERT_FILE", "certs/crt.pem")
	keyFile := env.GetEnv("TLS_KEY_FILE", "certs/key.pem")
	memory := env.GetEnv("IN_MEMORY", "false") == "true"

	var db *badger.DB
	var err error
	var options badger.Options

	if memory {
		logger.Println("running in memory mode")

		options = badger.DefaultOptions("").WithLogger(nil).WithInMemory(true)
	} else {
		options = badger.DefaultOptions(path.Join(data, "fsm")).WithLogger(nil).WithInMemory(false)
	}

	db, err = badger.Open(options)
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

	transport.RunHTTPServer(eventstore, logger)

	promServer := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	promServer.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	go func() {
		logger.Printf("Prometheus server listening on %s", promPort)

		if tlsEnabled {
			if err := promServer.ListenTLS(":"+promPort, certFile, keyFile); err != nil {
				logger.Fatal(err)
			}
		} else {
			if err := promServer.Listen(":" + promPort); err != nil {
				logger.Fatal(err)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	logger.Println("eventflowDB is shutting down...")

	db.Close()
	grpcServer.GracefulStop()
}

func main() {
	server()
}
