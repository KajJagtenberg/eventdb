package client

func (c *Client) Version() (string, error) {
	return c.r.Do("VERSION").String()
}
