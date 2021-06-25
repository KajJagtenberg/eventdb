package main

import (
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/dgraph-io/badger/v3"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/eventflowdb/eventflowdb/api"
	"github.com/eventflowdb/eventflowdb/env"
	"github.com/eventflowdb/eventflowdb/store"
	"github.com/eventflowdb/eventflowdb/transport"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	grpcPort := env.GetEnv("GRPC_PORT", "6543")
	httpPort := env.GetEnv("HTTP_PORT", "16543")
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

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		logger.Fatal(err)
	}

	grpcOptions := []grpc.ServerOption{}

	if tlsEnabled {
		logger.Println("tls is enabled")

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			logger.Fatal(err)
		}

		config := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.NoClientCert,
		}

		grpcOptions = append(grpcOptions, grpc.Creds(credentials.NewTLS(config)))
	}

	grpcServer := grpc.NewServer(grpcOptions...)

	api.RegisterEventStoreServer(grpcServer, transport.NewEventStore(eventstore))

	go func() {
		logger.Printf("gRPC server listening on %s", grpcPort)

		grpcServer.Serve(lis)
	}()

	httpServer := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	httpServer.Post("/api/v1", transport.HTTPHandler(eventstore, logger))

	go func() {
		logger.Printf("HTTP server listening on %s", httpPort)

		if tlsEnabled {
			if err := httpServer.ListenTLS(":"+httpPort, certFile, keyFile); err != nil {
				logger.Fatal(err)
			}
		} else {
			if err := httpServer.Listen(":" + httpPort); err != nil {
				logger.Fatal(err)
			}
		}
	}()

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
