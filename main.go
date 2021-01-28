package main

import (
	"bytes"
	"eventdb/env"
	"io/ioutil"
	"log"
	"net/http"
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

func LoadBabel(vm *goja.Runtime) error {
	r, err := http.Get("https://cdnjs.cloudflare.com/ajax/libs/babel-standalone/6.26.0/babel.min.js")
	if err != nil {
		return err
	}
	src, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	_, err = vm.RunString(string(src))
	return err
}

func Transpile(code string) (string, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	if err := LoadBabel(vm); err != nil {
		return "", err
	}

	vm.Set("input", code)

	if _, err := vm.RunString(`var output = Babel.transform(input, {presets: ["es2015"]}).code;`); err != nil {
		return "", err
	}
	return vm.Get("output").String(), nil
}

func sandbox() {
	file, err := os.OpenFile("projections/index.js", os.O_RDONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	src, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	transpiled, err := Transpile(string(src))
	if err != nil {
		log.Fatal(err)
	}

	vm := goja.New()
	vm.Set("log", log.Println)
	if _, err := vm.RunString(transpiled); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// server()

	sandbox()
}
