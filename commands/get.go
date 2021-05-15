package commands

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
)

const (
	CMD_GET       = "get"
	CMD_GET_SHORT = "g"
)

type GetRequest struct {
	Stream  uuid.UUID `json:"stream"`
	Version uint32    `json:"version"`
	Limit   uint32    `json:"limit"`
}

type GetResponse struct {
	Events []store.Event `json:"events"`
}

func GetHandler(store store.EventStore) CommandHandler {
	return func(cmd Command) (interface{}, error) {
		if cmd.Args == nil {
			return nil, ErrInsufficientArguments
		}

		var req GetRequest
		if err := json.Unmarshal(cmd.Args, &req); err != nil {
			return nil, err
		}

		events, err := store.Get(req.Stream, req.Version, req.Limit)
		if err != nil {
			return nil, err
		}

		return GetResponse{events}, nil
	}
}
