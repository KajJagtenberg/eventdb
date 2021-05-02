package client

func (c *Client) Ping() (string, error) {
	return c.r.Ping().Result()
}
