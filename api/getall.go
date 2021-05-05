package api

import (
	"encoding/json"

	"github.com/KajJagtenberg/eventflowdb/store"
	"github.com/oklog/ulid"
)

type GetAllRequest struct {
	Offset ulid.ULID `json:"stream"`
	Limit  uint32    `json:"limit"`
}

func GetAll(s store.Store, c *Ctx) error {
	if len(c.Args) == 0 {
		return ErrInsufficientArguments
	}

	var req GetAllRequest
	if err := json.Unmarshal(c.Args[0], &req); err != nil {
		return err
	}

	events, err := s.GetAll(req.Offset, req.Limit)
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
