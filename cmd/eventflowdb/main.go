package main

import (
	"os"
	"os/signal"
	"syscall"

	"log"

	"github.com/joho/godotenv"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/env"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/tidwall/redcon"
	"go.etcd.io/bbolt"

	_ "embed"
)

var (
	stateLocation = env.GetEnv("STATE_LOCATION", "data/state.dat")
	respAddr      = env.GetEnv("RESP_ADDR", ":6543")
	// httpAddr      = env.GetEnv("HTTP_ADDR", ":16543")
)

func main() {
	godotenv.Load()

	log.Println("Initializing store")

	db, err := bbolt.Open(stateLocation, 0666, bbolt.DefaultOptions)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	store, err := store.NewBoltStore(db)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// log.Println("Initializing HTTP server")

	// app := fiber.New(fiber.Config{
	// 	DisableStartupMessage: true,
	// })
	// app.Use(helmet.New())
	// app.Use(cors.New())

	// prom := adaptor.HTTPHandler(promhttp.Handler())

	// app.Get("/metrics", prom)

	// go func() {
	// 	log.Printf("HTTP server listening on %v", httpAddr)

	// 	if err := app.Listen(httpAddr); err != nil {
	// 		log.Fatalf("Failed to listen: %v", err)
	// 	}
	// }()

	log.Println("Initializing RESP server")

	resp := api.NewResp(store)

	go func() {
		log.Printf("RESP API listening on %s", respAddr)

		if err := redcon.ListenAndServe(respAddr, resp.CommandHandler, resp.AcceptHandler, resp.ErrorHandler); err != nil {
			log.Fatalf("Failed to run RESP API: %v", err)
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	log.Println("Shutting down...")
}
