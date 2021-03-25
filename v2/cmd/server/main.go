package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/kajjagtenberg/eventflowdb/cluster"
	"github.com/kajjagtenberg/eventflowdb/env"
	"github.com/kajjagtenberg/eventflowdb/graph/generated"
	"github.com/kajjagtenberg/eventflowdb/graph/resolvers"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.etcd.io/bbolt"
	"google.golang.org/grpc"
)

var (
	localID   = env.GetEnv("RAFT_LOCAL_ID", "")
	bindAddr  = env.GetEnv("RAFT_BIND_ADDR", ":6542")
	advrAddr  = env.GetEnv("RAFT_ADVR_ADDR", bindAddr)
	bootstrap = env.GetEnv("RAFT_BOOTSTRAP", "false") == "true"
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

	///////////////
	//  Storage  //
	///////////////

	log.Printf("Initializing Storage service at: %v", eventsFile)

	db, err := bbolt.Open(eventsFile, 0777, nil)
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

	if len(bindAddr) == 0 {
		log.Fatal("RAFT_BIND_ADDR cannot be empty")
	}

	fsm, err := cluster.NewFSM()
	if err != nil {
		log.Fatal(err)
	}

	raftServer, err := cluster.NewRaftServer(localID, bindAddr, advrAddr, fsm, bootstrap)
	if err != nil {
		log.Fatal(err)
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

	store.RegisterStreamsServer(grpcSrv, store.NewEventStoreService(storage))

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

	httpSrv.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/")))
	httpSrv.Post("/", adaptor.HTTPHandler(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{
		Storage: storage,
		Start:   time.Now(),
	}}))))
	httpSrv.Get("/backup", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/octet-stream")
		c.Set("Content-Disposition", "attachment;filename=backup.db")

		return storage.Backup(c.Response().BodyWriter())
	})

	//

	httpSrv.Get("/cluster", func(c *fiber.Ctx) error {
		return c.JSON(raftServer.Stats())
	})

	//

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
