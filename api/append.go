package api

import (
	"errors"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
)

type AppendRequest struct {
	Stream  uuid.UUID         `json:"stream"`
	Version uint32            `json:"version"`
	Events  []store.EventData `json:"events"`
}

type AppendResponse struct {
	Events []store.Event `json:"events"`
}

func AppendHandler(store store.Store) Handler {
	return func(c *Ctx) error {
		return errors.New("Not implemented")
	}
}
