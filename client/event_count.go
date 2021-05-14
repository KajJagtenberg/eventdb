package client

import "github.com/kajjagtenberg/eventflowdb/commands"

func (c *Client) EventCount() (int64, error) {
	return c.r.Do(commands.CMD_EVENT_COUNT_SHORT).Int64()
}

func (c *Client) EventCountEstimate() (int64, error) {
	return c.r.Do(commands.CMD_EVENT_COUNT_EST_SHORT).Int64()
}
