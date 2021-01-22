package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"eventdb/store"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func LoadFromStream(eventstore *store.Store) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		streamParam := vars["stream"]

		if len(streamParam) == 0 {
			http.Error(rw, "Stream cannot be empty", http.StatusBadRequest)
			return
		}

		stream, err := uuid.Parse(streamParam)
		if err != nil {
			http.Error(rw, "Stream must be an UUID v4", http.StatusBadRequest)
			return
		}

		versionQuery := r.FormValue("version")
		limitQuery := r.FormValue("limit")

		version, _ := strconv.Atoi(versionQuery)
		limit, _ := strconv.Atoi(limitQuery)

		if version < 0 {
			http.Error(rw, "Version cannot be negative", http.StatusBadRequest)
			return
		}

		if limit < 0 {
			http.Error(rw, "Limit cannot be negative", http.StatusBadRequest)
			return
		}

		events, err := eventstore.LoadFromStream(stream, version, limit)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		if len(events) == 0 {
			http.Error(rw, "Not Found", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(rw).Encode(events); err != nil {
			log.Println(err)
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func AppendToStream(eventstore *store.Store) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		streamParam := vars["stream"]
		versionParam := vars["version"]

		if len(streamParam) == 0 {
			http.Error(rw, "Stream cannot be empty", http.StatusBadRequest)
			return
		}

		stream, err := uuid.Parse(streamParam)
		if err != nil {
			http.Error(rw, "Stream must be an UUID v4", http.StatusBadRequest)
			return
		}

		version, _ := strconv.Atoi(versionParam)

		if version < 0 {
			http.Error(rw, "Version cannot be negative", http.StatusBadRequest)
			return
		}

		var events []store.AppendEvent

		if err := json.NewDecoder(r.Body).Decode(&events); err != nil {
			http.Error(rw, "Unable to decode body", http.StatusBadRequest)
			return
		}

		if len(events) == 0 {
			http.Error(rw, "Empty events", http.StatusBadRequest)
			return
		}

		validate := validator.New()

		for _, event := range events {
			if err := validate.Struct(event); err != nil {
				validationErrors := err.(validator.ValidationErrors)

				http.Error(rw, validationErrors.Error(), http.StatusBadRequest)
				return
			}
		}

		if err := eventstore.AppendToStream(stream, version, events); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		io.WriteString(rw, "Events added")
	}
}

func GetStreams(eventstore *store.Store) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		streams, err := eventstore.GetStreams(0, 1000)

		if err != nil {
			log.Println(err)
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(rw).Encode(streams); err != nil {
			log.Println(err)
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
