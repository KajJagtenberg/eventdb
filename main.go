package main

import (
	"encoding/json"
	"eventdb/env"
	"eventdb/projections"
	"io/ioutil"
	"log"
	"os"
	"time"

	"eventdb/handlers"
	"eventdb/store"

	"github.com/dop251/goja"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
	"github.com/oklog/ulid"
	"go.etcd.io/bbolt"
)

func setupMiddlewares(app *fiber.App) {
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{}))
	app.Use(etag.New())
}

func setupRoutes(app *fiber.App, eventstore *store.Store) {

	v1 := app.Group("/api/v1")
	v1.Use(logger.New(logger.Config{
		TimeZone: env.GetEnv("TZ", "UTC"),
	}))

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
	check(err)
	defer db.Close()

	eventstore, err := store.NewStore(db)
	check(err)

	log.Println("EventDB initializing API layer")

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	setupMiddlewares(app)
	setupRoutes(app, eventstore)

	go func() {
		compiler, err := projections.NewCompiler()
		check(err)

		code, err := LoadFileAsString("projections/index.js")
		check(err)

		code, err = compiler.Compile(code)
		check(err)

		var handlers map[string]func(event interface{})

		vm := goja.New()
		vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
		vm.Set("when", func(_handlers map[string]func(event interface{})) {
			handlers = _handlers
		})

		_, err = vm.RunString(code)
		check(err)

		collections := map[string]map[string]interface{}{}

		app.Get("/projections", compress.New(), func(c *fiber.Ctx) error {
			return c.JSON(collections)
		})

		vm.Set("set", func(collection string, id string, state map[string]interface{}) {
			_collection := collections[collection]

			if _collection == nil {
				_collection = map[string]interface{}{}
			}

			if state["version"] == nil {
				state["version"] = 1
			} else {
				state["version"] = state["version"].(int) + 1
			}

			_collection[id] = state

			collections[collection] = _collection
		})
		vm.Set("get", func(collection string, id string) interface{} {
			_collection := collections[collection]

			if _collection == nil {
				return map[string]interface {
				}{"version": 0}
			}

			state := _collection[id]

			if state == nil {
				return map[string]interface {
				}{"version": 0}
			} else {
				return state
			}
		})

		checkpoint := ulid.ULID{}

		for {
			events, err := eventstore.Subscribe(checkpoint, 10)
			check(err)

			if len(events) == 0 {
				time.Sleep(time.Second)
				continue
			}

			for _, event := range events {
				var data map[string]interface{}

				check(json.Unmarshal(event.Data, &data))

				arg := struct { // TODO: Add metadata (does not currently work with the way it's stored in the database)
					ID            string                 `json:"id"`
					Stream        string                 `json:"stream"`
					Version       int                    `json:"version"`
					Type          string                 `json:"type"`
					Data          map[string]interface{} `json:"data"`
					CausationID   string                 `json:"causation_id"`
					CorrelationID string                 `json:"correlation_id"`
				}{
					ID:            event.ID.String(),
					Stream:        event.Stream.String(),
					Version:       event.Version,
					Type:          event.Type,
					Data:          data,
					CausationID:   event.CausationID,
					CorrelationID: event.CorrelationID,
				}

				handler := handlers[event.Type]

				if handler != nil {
					handler(arg)
				}

				if handler := handlers["$any"]; handler != nil {
					handler(arg)
				}

				checkpoint = event.ID
			}
		}
	}()

	log.Println("EventDB initializing projection module")

	addr := env.GetEnv("LISTENING_ADDRESS", ":6543")

	log.Println("EventDB API layer ready to accept requests")

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
