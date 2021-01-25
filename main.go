package main

import (
	"eventdb/env"
	"log"

	"eventdb/handlers"
	"eventdb/store"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/helmet/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.etcd.io/bbolt"
)

var (
	requestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_requests_total",
	})
)

func setupRoutes(app *fiber.App, eventstore *store.Store) {
	app.Use(helmet.New(helmet.Config{}))
	app.Use(cors.New())
	app.Use(etag.New())

	app.Static("/", "./webui/out")

	v1 := app.Group("/api/v1")

	v1.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	app.Use(func(c *fiber.Ctx) error {
		requestCounter.Inc()

		return c.Next()
	})

	v1.Get("/", handlers.Home(eventstore))
	v1.Get("/streams", handlers.GetStreams(eventstore))
	v1.Get("/streams/all", handlers.Subscribe(eventstore))
	v1.Get("/streams/:stream", handlers.LoadFromStream(eventstore))
	v1.Post("/streams/:stream/:version", handlers.AppendToStream(eventstore))
	v1.Get("/count", handlers.GetEventCount(eventstore))
	v1.Get("/backup", handlers.Backup(eventstore))
}

func main() {
	log.Println("EventDB initializing storage layer")

	file := env.GetEnv("DATABASE_FILE", "data.bolt")

	db, err := bbolt.Open(file, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	eventstore, err := store.NewStore(db)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("EventDB initializing API layer")

	app := fiber.New(fiber.Config{})

	setupRoutes(app, eventstore)

	addr := env.GetEnv("LISTENING_ADDRESS", ":6543")

	app.Listen(addr)
}
