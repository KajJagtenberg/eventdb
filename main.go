package main

import (
	"eventdb/env"
	"log"

	"eventdb/handlers"
	"eventdb/store"

	"github.com/gofiber/fiber/v2"
	"go.etcd.io/bbolt"
)

func setupRoutes(app *fiber.App, eventstore *store.Store) {
	app.Get("/", handlers.Home(eventstore))
	app.Get("/streams", handlers.GetStreams(eventstore))
	app.Get("/streams/:stream", handlers.LoadFromStream(eventstore))
	app.Post("/streams/:stream/:version", handlers.AppendToStream(eventstore))
	app.Get("/count", handlers.GetEventCount(eventstore))
	app.Get("/backup", handlers.Backup(eventstore))
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

	app := fiber.New()

	setupRoutes(app, eventstore)

	addr := env.GetEnv("LISTENING_ADDRESS", ":5555")

	app.Listen(addr)
}
