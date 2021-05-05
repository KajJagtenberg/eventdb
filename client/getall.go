package client

import (
	"encoding/json"

	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/oklog/ulid"
)

func (c *Client) GetAll(offset ulid.ULID, limit uint32) ([]store.Event, error) {
	req := api.GetAllRequest{
		Offset: offset,
		Limit:  limit,
	}

	args, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := c.r.Do("GETALL", args).Result()
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
