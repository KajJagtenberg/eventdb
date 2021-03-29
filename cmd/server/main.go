package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/helmet/v2"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/cluster"
	"github.com/kajjagtenberg/eventflowdb/env"
	"github.com/kajjagtenberg/eventflowdb/graph/generated"
	"github.com/kajjagtenberg/eventflowdb/graph/resolvers"
	"github.com/kajjagtenberg/eventflowdb/persistence"
	"go.etcd.io/bbolt"
)

var (
	localID       = env.GetEnv("RAFT_LOCAL_ID", "main")
	bindAddr      = env.GetEnv("RAFT_BIND_ADDR", "127.0.0.1:6542")
	advrAddr      = env.GetEnv("RAFT_ADVR_ADDR", bindAddr)
	bootstrap     = env.GetEnv("RAFT_BOOTSTRAP", "false") == "true"
	stateLocation = env.GetEnv("STATE_LOCATION", "data/state.dat")
	grpcAddr      = env.GetEnv("GRPC_ADDR", ":6543")
	graphqlAddr   = env.GetEnv("GRAPHQL_ADDR", ":16543")
)

func main() {
	db, err := bbolt.Open(stateLocation, 0666, bbolt.DefaultOptions)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	persistence, err := persistence.NewPersistence(db)
	if err != nil {
		log.Fatalf("Failed to create persistence: %v", err)
	}

	fsm, err := cluster.NewFSM(persistence)
	if err != nil {
		log.Fatalf("Failed to create FSM: %v", err)
	}

	raftServer, err := cluster.NewRaftServer(localID, bindAddr, advrAddr, fsm, true)
	if err != nil {
		log.Fatalf("Failed to create Raft: %v", err)
	}
	defer raftServer.Shutdown()

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	grpcServer := api.NewGRPCServer(raftServer, persistence)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers.NewResolver(raftServer)}))

	http.Handle("/query", srv)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(helmet.New())
	app.Get("/graphql", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/graphql")))
	app.Post("/graphql", adaptor.HTTPHandler(srv))

	go func() {
		if err := app.Listen(graphqlAddr); err != nil {
			log.Fatalf("Failed to serve Graphql: %v", err)
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
