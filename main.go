package main

import (
	"eventdb/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"eventdb/handlers"
	"eventdb/store"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("data").WithLogger(nil))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	eventstore := store.NewStore(db)

	count := 1000000

	start := time.Now()

	for i := 0; i < count; i++ {
		cause := uuid.New().String()

		event := store.AppendEvent{
			Type:          "person_added",
			CausationID:   cause,
			CorrelationID: cause,
			Data: struct {
				Name string
			}{
				Name: "Kaj Jagtenberg",
			},
		}

		stream := uuid.New()

		if err := eventstore.AppendToStream(stream, 0, []store.AppendEvent{event}); err != nil {
			log.Fatal(err)
		}
	}

	passed := time.Now().Sub(start)

	log.Println("Speed: ", count/int(passed.Seconds()))

	router := mux.NewRouter()
	router.Use(middleware.JSONMiddleWare())
	router.Use(middleware.TimerMiddleWare())
	router.HandleFunc("/streams/{stream}", handlers.LoadFromStream(eventstore)).Methods(http.MethodGet)
	router.HandleFunc("/streams/{stream}/{version}", handlers.AppendToStream(eventstore)).Methods(http.MethodPost)
	router.HandleFunc("/streams", handlers.GetStreams(eventstore)).Methods(http.MethodGet)

	server := http.Server{
		Addr:         ":5555", // TODO: Get from env var
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("EventDB HTTP API listening on %s\n", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-c
}
