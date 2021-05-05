package api

import (
	"encoding/json"

	"github.com/KajJagtenberg/eventflowdb/store"
	"github.com/google/uuid"
)

type AddRequest struct {
	Stream  uuid.UUID         `json:"stream"`
	Version uint32            `json:"version"`
	Events  []store.EventData `json:"events"`
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

	c.Conn.WriteArray(len(events))

	for _, event := range events {
		v, err := json.Marshal(&event)
		if err != nil {
			return err
		}

		c.Conn.WriteString(string(v))
	}

	return nil
}
