package api

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
)

type GetRequest struct {
	Stream  uuid.UUID `json:"stream"`
	Version uint32    `json:"version"`
	Limit   uint32    `json:"limit"`
}

func Get(s store.Store, c *Ctx) error {
	if len(c.Args) == 0 {
		return ErrInsufficientArguments
	}

	var req GetRequest
	if err := json.Unmarshal(c.Args, &req); err != nil {
		return err
	}

	events, err := s.Get(req.Stream, req.Version, req.Limit)
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
