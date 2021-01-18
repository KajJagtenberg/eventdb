package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	store := NewStore(db)

	router := mux.NewRouter()
	router.HandleFunc("/streams/{stream}", func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		stream := vars["stream"]

		if len(stream) == 0 {
			http.Error(rw, "Stream name cannot be empty", http.StatusBadRequest)
			return
		}

		version := 0
		limit := 0

		// Validate request

		events, err := store.LoadFromStream(stream, version, limit)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		rw.Header().Set("Content-Type", "application/json;charset=utf-8")

		if err := json.NewEncoder(rw).Encode(events); err != nil {
			log.Println(err)
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
			return
		}
	})

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

	// TODO: Add an networked API
}
