package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"eventdb/store"

	"github.com/gorilla/mux"
)

func LoadFromStream(eventstore *store.Store) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		stream := vars["stream"]

		if len(stream) == 0 {
			http.Error(rw, "Stream name cannot be empty", http.StatusBadRequest)
			return
		}

		versionQuery := r.FormValue("version")
		limitQuery := r.FormValue("limit")

		version, _ := strconv.Atoi(versionQuery)
		limit, _ := strconv.Atoi(limitQuery)

		// Validate request

		events, err := eventstore.LoadFromStream(stream, version, limit)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.NewEncoder(rw).Encode(events); err != nil {
			log.Println(err)
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
