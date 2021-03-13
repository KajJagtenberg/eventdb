package main

import (
	"eventflowdb/env"
	"eventflowdb/graph/generated"
	graph "eventflowdb/graph/resolvers"
	"log"
	"time"

	"eventflowdb/store"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
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
	log.Println("EventflowDB initializing storage layer")

	file := env.GetEnv("DATABASE_FILE", "events.db")

	db, err := bbolt.Open(file, 0600, nil)
	check(err)
	defer db.Close()

	eventstore, err := store.NewEventStore(db)
	check(err)

	log.Println("EventflowDB initializing GraphQL")

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	setupMiddlewares(app)

	if env.GetEnv("DISABLE_PLAYGROUND", "false") != "true" {
		app.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/")))
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		EventStore: eventstore,
		DB:         db,
		Startup:    time.Now(),
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
	server()
}
