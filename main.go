package main

import (
	"eventdb/env"
	"eventdb/projections"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"eventdb/handlers"
	"eventdb/store"

	"github.com/dop251/goja"
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

	log.Println("EventDB initializing projection module")

	go func() {
		compiler, err := projections.NewCompiler()
		check(err)

		sourceFile, err := os.OpenFile("projections/index.js", os.O_RDONLY, 0600)
		check(err)

		sourceCode, err := ioutil.ReadAll(sourceFile)
		sourceFile.Close()

		code, err := compiler.Compile(string(sourceCode))
		check(err)

		fmt.Println(code)

		program, err := goja.Compile("index", code, false)
		check(err)

		vm := goja.New()
		vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
		vm.Set("println", func(a ...interface{}) {
			fmt.Println(a...)
		})

		type WhenInput struct {
			Initial map[string]interface{} `json:"$init"`
		}

		output, err := vm.RunProgram(program)
		check(err)

		fmt.Println(output.Export())
	}()

	addr := env.GetEnv("LISTENING_ADDRESS", ":6543")

	log.Println("EventDB API layer ready to accept requests")

	app.Listen(addr)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	server()
}
