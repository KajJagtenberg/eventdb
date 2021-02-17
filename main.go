package main

import (
	"eventflowdb/env"
	"io/ioutil"
	"log"
	"os"
	"time"

	"eventflowdb/handlers"
	"eventflowdb/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
	"go.etcd.io/bbolt"
)

func setupMiddlewares(app *fiber.App) {
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{}))
	app.Use(etag.New())
}

func setupRoutes(app *fiber.App, eventstore *store.EventStore) {
	app.Get("/", handlers.Home(eventstore))

	v1 := app.Group("/api/v1")
	v1.Use(logger.New(logger.Config{
		TimeZone: env.GetEnv("TZ", "UTC"),
	}))

	v1.Get("/streams", handlers.GetStreams(eventstore))
	v1.Get("/streams/all", handlers.Subscribe(eventstore))
	v1.Get("/streams/:stream", handlers.LoadFromStream(eventstore))
	v1.Post("/streams/:stream/:version", handlers.AppendToStream(eventstore))
	v1.Get("/events/:id", etag.New(), cache.New(cache.Config{
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}), handlers.GetEventByID(eventstore))
	v1.Get("/count", handlers.GetEventCount(eventstore))
	v1.Get("/backup", handlers.Backup(eventstore))
}

func server() {
	log.Println("EventDB initializing storage layer")

	file := env.GetEnv("DATABASE_FILE", "events.db")

	db, err := bbolt.Open(file, 0600, nil)
	check(err)
	defer db.Close()

	eventstore, err := store.NewEventStore(db)
	check(err)

	log.Println("EventDB initializing API layer")

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	setupMiddlewares(app)
	setupRoutes(app, eventstore)

	log.Println("EventDB initializing projection module")

	addr := env.GetEnv("LISTENING_ADDRESS", ":6543")

	log.Printf("EventDB API layer ready to accept requests on %s\n", addr)

	app.Listen(addr)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/* */

func LoadFileAsString(file string) (string, error) {
	fin, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		return "", err
	}
	defer fin.Close()

	src, err := ioutil.ReadAll(fin)
	if err != nil {
		return "", err
	}

	return string(src), nil
}

func main() {
	server()
}
