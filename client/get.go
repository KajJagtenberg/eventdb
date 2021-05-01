package client

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/store"
)

func (c *Client) Get(stream uuid.UUID, version uint32, limit uint32) ([]store.Event, error) {
	req := api.GetRequest{
		Stream:  stream,
		Version: version,
		Limit:   limit,
	}

	args, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := c.r.Do("GET", args).Result()
	if err != nil {
		return nil, err
	}

	var events []store.Event

	for _, entry := range response.([]interface{}) {
		var event store.Event
		if err := json.Unmarshal([]byte(entry.(string)), &event); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}
