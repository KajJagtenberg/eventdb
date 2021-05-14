package client

import "github.com/kajjagtenberg/eventflowdb/commands"

func (c *Client) Version() (string, error) {
	return c.r.Do(commands.CMD_VERSION_SHORT).String()
}
