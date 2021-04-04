package main

import (
	"os"
	"os/signal"
	"syscall"

	"log"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"
	"github.com/kajjagtenberg/eventflowdb/env"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.etcd.io/bbolt"
)

var (
	// localID       = env.GetEnv("RAFT_LOCAL_ID", "main")
	// bindAddr      = env.GetEnv("RAFT_BIND_ADDR", "127.0.0.1:6542")
	// advrAddr      = env.GetEnv("RAFT_ADVR_ADDR", bindAddr)
	// bootstrap     = env.GetEnv("RAFT_BOOTSTRAP", "false") == "true"
	stateLocation = env.GetEnv("STATE_LOCATION", "data/state.dat")
	grpcAddr      = env.GetEnv("GRPC_ADDR", ":6543")
	httpAddr      = env.GetEnv("HTTP_ADDR", ":16543")
	// graphqlAddr   = env.GetEnv("GRAPHQL_ADDR", ":16543")
)

func main() {
	db, err := bbolt.Open(stateLocation, 0666, bbolt.DefaultOptions)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	store, err := store.NewBoltStore(db)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	log.Println(store.Size()) // TODO: Remove when no longer necessary

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

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
