package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base32"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/dgraph-io/badger/v3"
	"github.com/kajjagtenberg/eventflowdb/env"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/kajjagtenberg/eventflowdb/transport"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	data       = env.GetEnv("DATA", "data")
	port       = env.GetEnv("PORT", "6543")
	password   = env.GetEnv("PASSWORD", "")
	noPassword = env.GetEnv("NO_PASSWORD", "false") == "true"
	tlsEnabled = env.GetEnv("TLS_ENABLED", "false") == "true"
	certFile   = env.GetEnv("TLS_CERT_FILE", "certs/cert.pem")
	keyFile    = env.GetEnv("TLS_KEY_FILE", "certs/key.pem")
	production = env.GetEnv("ENVIRONMENT", "development") == "production"

	log = logrus.New()
)

func check(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {
	log.SetFormatter(&logrus.JSONFormatter{})

	if production {
		if !noPassword && len(password) == 0 {
			passwordData := make([]byte, 20)
			_, err := rand.Read(passwordData)
			check(err, "failed to generate password")

			password = base32.StdEncoding.EncodeToString(passwordData)

			log.Printf("generated a password since none was given: %s", password)
		}
	}

	log.Println("initializing store")

	db, err := badger.Open(badger.DefaultOptions(path.Join(data)).WithLogger(log))
	check(err, "failed to open database")
	defer db.Close()

	eventstore, err := store.NewBadgerEventStore(store.BadgerStoreOptions{
		DB:             db,
		EstimateCounts: true,
	})
	check(err, "failed to create store")
	defer eventstore.Close()

	log.Println("initializing RESP server")

	var tlsConfig *tls.Config

	if tlsEnabled {
		crt, err := tls.LoadX509KeyPair(certFile, keyFile)
		check(err, "failed to load certificate")

		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{crt},
		}
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(transport.AuthenticationInterceptor))

	transport.RegisterEventStoreServiceServer(grpcServer, transport.NewEventStoreService(eventstore))

	go func() {
		if tlsEnabled {
			l, err := tls.Listen("tcp", ":"+port, tlsConfig)
			check(err, "failed to create listener")

			log.Printf("gRPC API listening on %s over TLS", port)

			check(grpcServer.Serve(l), "failed to run gRPC API over TLS")
		} else {
			l, err := net.Listen("tcp", ":"+port)
			check(err, "failed to create listener")

			log.Printf("gRPC API listening on %s", port)

			check(grpcServer.Serve(l), "failed to run gRPC API")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	<-c

	log.Println("eventflowDB is shutting down...")
}
