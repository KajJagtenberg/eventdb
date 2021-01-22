package handlers

import (
	"encoding/json"
	"eventdb/store"
	"net/http"
)

func GetEventCOunt(eventstore *store.Store) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		count, err := eventstore.GetEventCount()
		if err != nil {
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(rw).Encode(struct {
			Count int `json:"count"`
		}{count}); err != nil {
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
		}
	}
}
