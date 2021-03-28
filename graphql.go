package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/helmet/v2"
	"github.com/kajjagtenberg/eventflowdb/graph"
	"github.com/kajjagtenberg/eventflowdb/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/query", srv)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(helmet.New())
	app.Get("/graphql", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/graphql")))
	app.Post("/graphql", adaptor.HTTPHandler(srv))

	log.Printf("Connect to http://localhost:%s/graphql for GraphQL playground", port)
	log.Fatal(app.Listen(":" + port))
}
