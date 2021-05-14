package client

import "github.com/kajjagtenberg/eventflowdb/commands"

func (c *Client) Size() (int64, string, error) {
	res, err := c.r.Do(commands.CMD_SIZE_SHORT).Result()
	if err != nil {
		return 0, "", err
	}

	entries := res.([]interface{})

	size, ok := entries[0].(int64)
	if !ok {
		return 0, "", ErrInvalidResponse
	}

	human, ok := entries[1].(string)
	if !ok {
		return 0, "", ErrInvalidResponse
	}

	return size, human, nil
}
