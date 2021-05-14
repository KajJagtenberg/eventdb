package client

import (
	"encoding/base64"
	"encoding/json"

	"github.com/kajjagtenberg/eventflowdb/commands"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/oklog/ulid"
)

func (c *Client) GetAll(offset ulid.ULID, limit uint32) ([]store.Event, error) {
	req := commands.GetAllRequest{
		Offset: offset,
		Limit:  limit,
	}

	cmd, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := c.r.Do(commands.CMD_GET_ALL_SHORT, base64.StdEncoding.EncodeToString(cmd)).Result()
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
