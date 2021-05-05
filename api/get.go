package api

import (
	"encoding/json"

	"github.com/KajJagtenberg/eventflowdb/store"
	"github.com/google/uuid"
)

type GetRequest struct {
	Stream  uuid.UUID `json:"stream"`
	Version uint32    `json:"version"`
	Limit   uint32    `json:"limit"`
}

type GetResponse struct {
	Events []store.Event `json:"events"`
}

func Get(s store.Store, c *Ctx) error {
	if len(c.Args) == 0 {
		return ErrInsufficientArguments
	}

	var req GetRequest
	if err := json.Unmarshal(c.Args[0], &req); err != nil {
		return err
	}

	events, err := s.Get(req.Stream, req.Version, req.Limit)
	if err != nil {
		return err
	}

	res := GetResponse{
		Events: events,
	}

	result, err := json.Marshal(&res)
	if err != nil {
		return err
	}

	c.Conn.WriteString(string(result))

	return nil
}
