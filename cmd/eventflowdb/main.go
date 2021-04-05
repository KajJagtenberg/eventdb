package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/env"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.etcd.io/bbolt"
	"google.golang.org/grpc"
)

var (
	stateLocation = env.GetEnv("STATE_LOCATION", "data/state.dat")
	grpcAddr      = env.GetEnv("GRPC_ADDR", ":6543")
	httpAddr      = env.GetEnv("HTTP_ADDR", ":16543")
)

func main() {
	log.Println("Initializing store")

	db, err := bbolt.Open(stateLocation, 0666, bbolt.DefaultOptions)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	store, err := store.NewBoltStore(db)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	log.Println("Initializing HTTP server")

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(helmet.New())
	app.Use(cors.New())

	prom := adaptor.HTTPHandler(promhttp.Handler())

	app.Get("/metrics", prom)

	go func() {
		log.Printf("HTTP server listening on %v", httpAddr)

		if err := app.Listen(httpAddr); err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
	}()

	log.Println("Initializing gRPC server")

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen on socket: %v", err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()

	api.RegisterStreamServiceServer(grpcServer, api.NewStreamService(store))
	api.RegisterShellServiceServer(grpcServer, api.NewShellService(store))

	go func() {
		log.Printf("gRPC server listening on %v", grpcAddr)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to listen on gRPC server: %v", err)
		}
	}()

	go func() {
		for {
			streams, err := store.StreamCount()
			if err != nil {
				log.Fatalf("Failed to get stream count: %v", err)
			}

			events, err := store.EventCount()
			if err != nil {
				log.Fatalf("Failed to get event count: %v", err)
			}

			log.Printf("Streams: %d, Events: %d", streams, events)

			time.Sleep(time.Second)
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	log.Println("Shutting down...")
}
