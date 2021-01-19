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
	"github.com/gorilla/mux"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("data").WithLogger(nil))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// db.View(func(txn *badger.Txn) error {
	// 	opts := badger.DefaultIteratorOptions
	// 	opts.PrefetchSize = 10
	// 	it := txn.NewIterator(opts)
	// 	defer it.Close()
	// 	for it.Rewind(); it.Valid(); it.Next() {
	// 		log.Println(it.Item().Key())
	// 	}

	// 	return nil
	// })

	eventstore := store.NewStore(db)

	log.Println(db.Size())

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
