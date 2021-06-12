package main

import (
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/dgraph-io/badger/v3"
	"github.com/joho/godotenv"
	"github.com/kajjagtenberg/eventflowdb/env"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/kajjagtenberg/eventflowdb/transport"
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

func main() {
	server()
}
