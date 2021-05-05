package client

func (c *Client) Uptime() (string, error) {
	return c.r.Do("UPTIME").String()
}
