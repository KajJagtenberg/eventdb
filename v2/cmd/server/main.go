package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hashicorp/memberlist"
	"github.com/joho/godotenv"
	"github.com/kajjagtenberg/eventflowdb/env"
	"github.com/kajjagtenberg/eventflowdb/graph/generated"
	"github.com/kajjagtenberg/eventflowdb/graph/resolvers"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.etcd.io/bbolt"
	"google.golang.org/grpc"
)

func main() {
	/////////////
	//  Hello  //
	/////////////

	log.Println("Hello, world!")

	//////////////
	//  Config  //
	//////////////

	godotenv.Load()

	grpcAddr := env.GetEnv("GRPC_LISTENER", ":6543")
	httpAddr := env.GetEnv("HTTP_LISTENER", ":16543")
	eventsFile := env.GetEnv("EVENTS_FILE", "events.db")
	existingNodes := env.GetEnv("EXISTING_NODES", "")

	///////////////
	//  Storage  //
	///////////////

	log.Println("Initializing Storage service")

	db, err := bbolt.Open(eventsFile, 0666, nil)
	if err != nil {
		log.Fatalf("Failed to initialize Storage service: %v", err)
	}
	defer db.Close()

	storage, err := store.NewStorage(db)
	if err != nil {
		log.Fatalf("Failed to initialize Storage service: %v", err)
	}

	///////////////
	//  Cluster  //
	///////////////

	log.Println("Setting up a cluster")

	conf := memberlist.DefaultLocalConfig()

	cluster, err := memberlist.Create(conf)
	if err != nil {
		log.Fatalf("Failed to create cluster: %v", err)
	}
	defer cluster.Leave(time.Second * 5)

	var existing []string

	if len(existingNodes) > 0 {
		existing = strings.Split(existingNodes, ",")
	}

	if joined, err := cluster.Join(existing); err != nil {
		log.Fatalf("Failed to join a cluster: %v", err)
	} else {
		log.Printf("Successfully joined a cluster with %d nodes", joined)
	}

	////////////
	//  gRPC  //
	////////////

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s", grpcAddr)
	}
	defer lis.Close()

	grpcSrv := grpc.NewServer()

	log.Println("Initializing gRPC services")

	store.RegisterEventStoreServer(grpcSrv, store.NewStoreService(storage))

	go func() {
		log.Printf("Starting gRPC server on %s", grpcAddr)

		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	////////////
	//  HTTP  //
	////////////

	log.Println("Initializing Prometheus metrics")

	httpSrv := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ReadTimeout:           time.Second * 10,
		WriteTimeout:          time.Second * 10,
		IdleTimeout:           time.Second * 10,
	})
	httpSrv.Use(cors.New())
	httpSrv.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	log.Println("Initializing GraphQL")

	httpSrv.Get("/graphql", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/graphql")))
	httpSrv.Post("/graphql", adaptor.HTTPHandler(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{
		Memberlist: cluster,
		Storage:    storage,
	}}))))

	go func() {
		log.Printf("Starting HTTP server on %s", httpAddr)

		httpSrv.Listen(httpAddr)
	}()

	////////////////
	//  Shutdown  //
	////////////////

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	<-sig

	log.Println("Stopping all services")

	grpcSrv.GracefulStop()
	// httpSrv.Shutdown()
	db.Close()

	log.Println("Stopped Storage")

	log.Println("Stopped all services")
}
