package main

import (
	"context"
	"eventdb/env"
	"eventdb/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"eventdb/handlers"
	"eventdb/store"

	"github.com/gorilla/mux"
	"go.etcd.io/bbolt"
)

func main() {
	log.Println("EventDB initializing storage layer")

	db, err := bbolt.Open("data.bolt", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	eventstore, err := store.NewStore(db)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("EventDB initializing API layer")

	router := mux.NewRouter()
	router.Use(middleware.JSONMiddleWare())
	router.HandleFunc("/", handlers.Home())
	router.HandleFunc("/streams/{stream}", handlers.LoadFromStream(eventstore)).Methods(http.MethodGet)
	router.HandleFunc("/streams/{stream}/{version}", handlers.AppendToStream(eventstore)).Methods(http.MethodPost)
	router.HandleFunc("/streams", handlers.GetStreams(eventstore)).Methods(http.MethodGet)
	router.HandleFunc("/count", handlers.GetEventCount(eventstore)).Methods(http.MethodGet)
	router.HandleFunc("/backup", handlers.Backup(eventstore)).Methods(http.MethodGet)

	server := http.Server{
		Addr:         env.GetEnv("LISTENING_ADDRESS", ":5555"),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	defer server.Shutdown(context.Background())

	go func() {
		log.Printf("EventDB HTTP API listening on %s\n", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	AwaitShutdown()
}

func AwaitShutdown() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-c
}
