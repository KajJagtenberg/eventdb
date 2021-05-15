package client

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/commands"
	"github.com/kajjagtenberg/eventflowdb/store"
)

func (c *Client) Add(stream uuid.UUID, version uint32, data []store.EventData) ([]store.Event, error) {
	req := commands.AddRequest{
		Stream:  stream,
		Version: version,
		Events:  data,
	}

	cmd, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	response, err := c.r.Do(commands.CMD_ADD_SHORT, cmd).Result()
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
