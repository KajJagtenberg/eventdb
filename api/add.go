package api

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
)

type AddRequest struct {
	Stream  uuid.UUID         `json:"stream"`
	Version uint32            `json:"version"`
	Events  []store.EventData `json:"events"`
}

type AddResponse struct {
	Events []store.Event `json:"events"`
}

func Add(s store.Store, c *Ctx) error {
	if len(c.Args) == 0 {
		return ErrInsufficientArguments
	}

	var req AddRequest
	if err := json.Unmarshal(c.Args[0], &req); err != nil {
		return err
	}

	events, err := s.Add(req.Stream, req.Version, req.Events)
	if err != nil {
		return err
	}

	res := AddResponse{
		Events: events,
	}

	result, err := json.Marshal(&res)
	if err != nil {
		return err
	}

	c.Conn.WriteString(string(result))

	return nil
}
