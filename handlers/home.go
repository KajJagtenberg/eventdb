package handlers

import (
	"encoding/json"
	"eventdb/constants"
	"log"
	"net/http"
)

func Home() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(rw).Encode(struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}{Name: constants.Name, Version: constants.Version}); err != nil {
			log.Println(err)
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
		}
	}
}
