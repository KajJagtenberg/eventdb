package commands

import (
	"encoding/json"

	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/kajjagtenberg/go-commando"
	"github.com/oklog/ulid"
)

const (
	CMD_GET_ALL       = "getall"
	CMD_GET_ALL_SHORT = "ga"
)

type GetAllRequest struct {
	Offset ulid.ULID `json:"offset"`
	Limit  uint32    `json:"limit"`
}

type GetAllResponse struct {
	Events []store.Event `json:"events"`
}

func GetAllHandler(store store.EventStore) commando.CommandHandler {
	return func(cmd commando.Command) (interface{}, error) {
		if cmd.Args == nil || len(cmd.Args) == 0 {
			return nil, commando.ErrInsufficientArguments
		}

		var req GetAllRequest
		if err := json.Unmarshal(cmd.Args, &req); err != nil {
			return nil, err
		}

		events, err := store.GetAll(req.Offset, req.Limit)
		if err != nil {
			return nil, err
		}

		return GetAllResponse{events}, nil
	}
}

func SetupGetAllHandler(dispatcher *commando.CommandDispatcher, store store.EventStore) {
	dispatcher.Register(CMD_GET_ALL, CMD_GET_ALL_SHORT, GetAllHandler(store))
}
