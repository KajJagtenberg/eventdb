package client

import "github.com/kajjagtenberg/eventflowdb/commands"

func (c *Client) Uptime() (string, error) {
	return c.r.Do(commands.CMD_UPTIME_SHORT).String()
}
