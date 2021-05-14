package client

import "github.com/kajjagtenberg/eventflowdb/commands"

func (c *Client) StreamCount() (int64, error) {
	return c.r.Do(commands.CMD_STREAM_COUNT_SHORT).Int64()
}

func (c *Client) StreamCountEstimate() (int64, error) {
	return c.r.Do(commands.CMD_STREAM_COUNT_EST_SHORT).Int64()
}
