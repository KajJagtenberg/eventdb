package main

import (
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/joho/godotenv"
	"github.com/kajjagtenberg/eventflowdb/api"
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
	maxPool            = 3
	tcpTimeout         = 10 * time.Second
	raftSnapShotRetain = 2
	raftLogCacheSize   = 512
)

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})

	godotenv.Load()
}

func server() {
	data := env.GetEnv("DATA", "data")
	raftPort := env.GetEnv("RAFT_PORT", "26543")
	followers := env.GetEnv("FOLLOWERS", "")
	grpcPort := env.GetEnv("PORT", "6543")
	tlsEnabled := env.GetEnv("TLS_ENABLED", "false") == "true"
	certFile := env.GetEnv("TLS_CERT_FILE", "certs/cert.pem")
	keyFile := env.GetEnv("TLS_KEY_FILE", "certs/key.pem")

	db, err := badger.Open(badger.DefaultOptions(path.Join(data, "fsm")).WithLogger(log))
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

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	nodeID := env.GetEnv("NODE_ID", hostname)

	raftConf := raft.DefaultConfig()
	raftConf.SnapshotThreshold = 1024
	raftConf.LocalID = raft.ServerID(nodeID)

	fsmStore := fsm.NewFSM(eventstore)

	store, err := raftboltdb.NewBoltStore(path.Join(data, "store"))
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	cacheStore, err := raft.NewLogCache(raftLogCacheSize, store)
	if err != nil {
		log.Fatal(err)
	}

	snapshotStore, err := raft.NewFileSnapshotStore(path.Join(data, "snapshots"), raftSnapShotRetain, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	raftBindAddr := nodeID + ":" + raftPort

	tcpAddr, err := net.ResolveTCPAddr("tcp", raftBindAddr)
	if err != nil {
		log.Fatal(err)
	}

	raftTransport, err := raft.NewTCPTransport(raftBindAddr, tcpAddr, maxPool, tcpTimeout, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	raftSrv, err := raft.NewRaft(raftConf, fsmStore, cacheStore, store, snapshotStore, raftTransport)
	if err != nil {
		log.Fatal(err)
	}
	defer raftSrv.Shutdown()

	var configuration raft.Configuration
	configuration.Servers = []raft.Server{
		{
			ID:      raft.ServerID(raftConf.LocalID),
			Address: raftTransport.LocalAddr(),
		},
	}

	if len(followers) > 0 {
		for _, follower := range strings.Split(followers, ",") {
			parts := strings.Split(follower, "@")

			configuration.Servers = append(configuration.Servers,
				raft.Server{
					ID:      raft.ServerID(parts[0]),
					Address: raft.ServerAddress(parts[1]),
				})
		}
	}

	raftSrv.BootstrapCluster(configuration)

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

	api.RegisterEventStoreServiceServer(grpcServer, transport.NewEventStoreService(raftSrv))

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

func main() {
	server()
}
