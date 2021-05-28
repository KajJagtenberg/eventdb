package commands

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/kajjagtenberg/go-commando"
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

func AddHandler(store store.EventStore) commando.CommandHandler {
	return func(cmd commando.Command) (interface{}, error) {
		if cmd.Args == nil || len(cmd.Args) == 0 {
			return nil, commando.ErrInsufficientArguments
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

func SetupAddHandler(dispatcher *commando.CommandDispatcher, store store.EventStore) {
	dispatcher.Register(CMD_ADD, CMD_ADD_SHORT, AddHandler(store))
}
