package handlers

import (
	"compress/gzip"
	"eventdb/store"
	"net/http"
)

func Backup(eventstore *store.Store) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/octet-stream")
		rw.Header().Set("Content-Disposition", "attachment;filename=eventdb.bak")

		if err := eventstore.Backup(gzip.NewWriter(rw)); err != nil {
			http.Error(rw, "Internals erver error", http.StatusInternalServerError)
		}
	}
}
