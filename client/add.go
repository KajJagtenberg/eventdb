package client

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/store"
)

func (c *Client) Add(stream uuid.UUID, version uint32, data []store.EventData) ([]store.Event, error) {
	req := api.AddRequest{
		Stream:  stream,
		Version: version,
		Events:  data,
	}

	cmd, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	response, err := c.r.Do("ADD", base64.StdEncoding.EncodeToString(cmd)).Result()
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
