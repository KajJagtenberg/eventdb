package main

import (
	"eventflowdb/env"
	"eventflowdb/graph"
	"eventflowdb/graph/generated"
	"log"
	"net/http"

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

	log.Println("EventflowDB initializing HTTP API layer")

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		EventStore: eventstore,
	}}))

	setupMiddlewares(app)

	app.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/query")))
	app.Post("/query", adaptor.HTTPHandler(srv))
	http.Handle("/query", srv)

	log.Println("EventflowDB initializing projection module")

	addr := env.GetEnv("LISTENING_ADDRESS", ":6543")

	log.Printf("EventflowDB HTTP API layer ready to accept requests on %s\n", addr)

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
