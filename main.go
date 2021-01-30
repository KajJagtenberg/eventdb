package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"eventdb/env"
	"io/ioutil"
	"log"
	"os"
	"time"

	"eventdb/handlers"
	"eventdb/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

func setupMiddlewares(app *fiber.App) {
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{}))
	app.Use(logger.New(logger.Config{
		TimeZone: env.GetEnv("TZ", "UTC"),
	}))
}

func setupRoutes(app *fiber.App, eventstore *store.Store) {

	v1 := app.Group("/api/v1")

	v1.Get("/", handlers.Home(eventstore))
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
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	eventstore, err := store.NewStore(db)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("EventDB initializing API layer")

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	setupMiddlewares(app)
	setupRoutes(app, eventstore)

	addr := env.GetEnv("LISTENING_ADDRESS", ":6543")

	log.Println("EventDB API layer ready to accept requests")

	app.Listen(addr)
}

func LoadFile(name string) *bytes.Buffer {
	file, err := os.OpenFile(name, os.O_RDONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	return bytes.NewBuffer(data)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func sandbox() {
	db, err := bbolt.Open("events.db", 0600, nil)
	check(err)
	defer db.Close()

	eventstore, err := store.NewStore(db)
	check(err)

	input, err := os.OpenFile("events.jsonld", os.O_RDONLY, 0600)
	check(err)
	defer input.Close()

	reader := bufio.NewReader(input)

	for line, _, err := reader.ReadLine(); err == nil; line, _, err = reader.ReadLine() {
		event := struct {
			Stream        uuid.UUID       `json:"stream"`
			Version       int             `json:"version"`
			Type          string          `json:"type"`
			Data          json.RawMessage `json:"data"`
			Timestamp     string          `json:"ts"`
			Metadata      struct{}        `json:"metadata"`
			CausationID   string          `json:"causation_id"`
			CorrelationID string          `json:"correlation_id"`
		}{}

		check(json.Unmarshal(line, &event))
		if err := eventstore.AppendToStream(event.Stream, event.Version, []store.AppendEvent{
			{
				Type:          event.Type,
				Data:          event.Data,
				Metadata:      event.Metadata,
				CausationID:   event.CausationID,
				CorrelationID: event.CorrelationID,
			},
		}); err != nil {
			log.Println(err)
		}
	}
}

func main() {
	// sandbox()
	server()
}
