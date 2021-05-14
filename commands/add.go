package commands

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
)

const (
	CMD_ADD       = "add"
	CMD_ADD_SHORT = "a"
)

type AddRequest struct {
	Stream  uuid.UUID         `json:"stream"`
	Version uint32            `json:"version"`
	Events  []store.EventData `json:"events"`
}

type AddResponse struct {
	Events []store.Event `json:"events"`
}

func AddHandler(store store.Store) CommandHandler {
	return func(cmd Command) (interface{}, error) {
		if cmd.Args == nil {
			return nil, ErrInsufficientArguments
		}

		var req AddRequest
		if err := json.Unmarshal(cmd.Args, &req); err != nil {
			return nil, err
		}

		events, err := store.Add(req.Stream, req.Version, req.Events)
		if err != nil {
			return nil, err
		}

		return AddResponse{events}, nil
	}
}
