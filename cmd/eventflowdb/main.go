package main

import (
	"context"
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/dgraph-io/badger/v3"
	clientv3 "go.etcd.io/etcd/client/v3"
	concurrency "go.etcd.io/etcd/client/v3/concurrency"

	"github.com/joho/godotenv"
	"github.com/kajjagtenberg/eventflowdb/api"
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
	grpcPort := env.GetEnv("PORT", "6543")
	tlsEnabled := env.GetEnv("TLS_ENABLED", "false") == "true"
	certFile := env.GetEnv("TLS_CERT_FILE", "certs/cert.pem")
	keyFile := env.GetEnv("TLS_KEY_FILE", "certs/key.pem")
	memory := env.GetEnv("IN_MEMORY", "false") == "true"

	var db *badger.DB
	var err error

	if memory {
		db, err = badger.Open(badger.DefaultOptions("").WithLogger(log).WithInMemory(true))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		db, err = badger.Open(badger.DefaultOptions(path.Join(data, "fsm")).WithLogger(log))
		if err != nil {
			log.Fatal(err)
		}
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
	defer grpcServer.Stop()

	api.RegisterEventStoreServiceServer(grpcServer, transport.NewEventStoreService(eventstore))

	go func() {
		log.Printf("gRPC server listening on %s", grpcPort)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	<-c

	log.Println("eventflowDB is shutting down...")
}

func etcd() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	id := env.GetEnv("NODE_ID", hostname)

	log.Println("node id:", id)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	log.Println("connected to etcd")

	session, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	ctx := context.Background()

	election := concurrency.NewElection(session, "/leader-election/")

	leader, err := election.Leader(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(leader.Kvs)

	log.Println("trying to elect")

	if err := election.Campaign(ctx, id); err != nil {
		log.Fatal(err)
	}

	log.Println("elected as leader")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	<-c
}

func main() {
	// server()
	etcd()
}
