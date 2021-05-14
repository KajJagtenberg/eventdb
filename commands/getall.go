package commands

import (
	"encoding/json"

	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/oklog/ulid"
)

type GetAllRequest struct {
	Offset ulid.ULID `json:"stream"`
	Limit  uint32    `json:"limit"`
}

type GetAllResponse struct {
	Events []store.Event `json:"events"`
}

func GetAllHandler(store store.Store) CommandHandler {
	return func(cmd Command) (interface{}, error) {
		if cmd.Args == nil {
			return nil, ErrInsufficientArguments
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
