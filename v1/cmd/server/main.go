package main

import (
	"eventflowdb/env"
	"eventflowdb/graph/generated"
	graph "eventflowdb/graph/resolvers"
	"eventflowdb/projections"
	"log"
	"time"

	"eventflowdb/store"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/helmet/v2"
	"go.etcd.io/bbolt"
)

func setupMiddlewares(app *fiber.App) {
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{}))
	app.Use(etag.New())
}

func server() {
	log.Println("Initializing Event Store")

	eventstoreFile := env.GetEnv("EVENT_STORE_FILE", "events.db")

	eventstoreDB, err := bbolt.Open(eventstoreFile, 0600, nil)
	check(err)
	defer eventstoreDB.Close()

	eventstore, err := store.NewEventStore(eventstoreDB)
	check(err)

	log.Println("Initializing Projection Engine")

	projectionEngineFile := env.GetEnv("PROJECTION_ENGINE_FILE", "projections.db")

	projectionEngineDB, err := bbolt.Open(projectionEngineFile, 0600, nil)
	check(err)
	defer projectionEngineDB.Close()

	projectionEngine, err := projections.NewProjectionEngine(projectionEngineDB)
	check(err)

	log.Println("Initializing GraphQL")

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	setupMiddlewares(app)

	if env.GetEnv("DISABLE_PLAYGROUND", "false") != "true" {
		app.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/")))
	}

	app.Get("/backup", compress.New(), func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/octet-stream")
		c.Set("Content-Disposition", "attachment;filename=backup.db")

		return eventstore.Backup(c.Response().BodyWriter())
	})

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		EventStore:       eventstore,
		ProjectionEngine: projectionEngine,
		Startup:          time.Now(),
	}}))

	app.Post("/", adaptor.HTTPHandler(srv))

	addr := env.GetEnv("LISTENING_ADDRESS", ":6543")

	log.Printf("EventflowDB GraphQL layer ready to accept requests on %s\n", addr)

	check(app.Listen(addr))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// sandbox()
	server()
}