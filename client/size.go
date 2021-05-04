package client

import "errors"

func (c *Client) Size() (int64, string, error) {
	res, err := c.r.Do("SIZE").Result()
	if err != nil {
		return 0, "", err
	}

	entries := res.([]interface{})

	size, ok := entries[0].(int64)
	if !ok {
		return 0, "", errors.New("invalid response")
	}

	human, ok := entries[1].(string)
	if !ok {
		return 0, "", errors.New("invalid response")
	}

	return size, human, nil
}
