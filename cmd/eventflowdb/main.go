package main

import (
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/dgraph-io/badger/v3"

	"github.com/eventflowdb/eventflowdb/api"
	"github.com/eventflowdb/eventflowdb/env"
	"github.com/eventflowdb/eventflowdb/store"
	"github.com/eventflowdb/eventflowdb/transport"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	log = logrus.New()
)

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})

	godotenv.Load()
}

func setupTLS(lis net.Listener, certFile, keyFile string) (net.Listener, error) {
	crt, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{crt},
	}

	lis = tls.NewListener(lis, config)
	if err != nil {
		return nil, err
	}

	return lis, nil
}

func server() {
	data := env.GetEnv("DATA", "data")
	grpcPort := env.GetEnv("PORT", "6543")
	tlsEnabled := env.GetEnv("TLS_ENABLED", "false") == "true"
	certFile := env.GetEnv("TLS_CERT_FILE", "certs/cert.pem")
	keyFile := env.GetEnv("TLS_KEY_FILE", "certs/key.pem")
	memory := env.GetEnv("IN_MEMORY", "false") == "true"

	var db *badger.DB
	var err error
	var options badger.Options

	if memory {
		log.Println("running in memory mode")

		options = badger.DefaultOptions("").WithLogger(nil).WithInMemory(true)
	} else {
		options = badger.DefaultOptions(path.Join(data, "fsm")).WithLogger(nil).WithInMemory(false)
	}

	db, err = badger.Open(options)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	eventstore, err := store.NewBadgerEventStore(store.BadgerStoreOptions{
		DB:             db,
		EstimateCounts: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer eventstore.Close()

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	if tlsEnabled {
		lis, err = setupTLS(lis, certFile, keyFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()

	api.RegisterEventStoreServer(grpcServer, transport.NewEventStore(eventstore))

	go func() {
		log.Printf("gRPC server listening on %s", grpcPort)

		grpcServer.Serve(lis)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	log.Println("eventflowDB is shutting down...")

	db.Close()
	grpcServer.GracefulStop()
}

func main() {
	server()
}
