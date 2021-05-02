package client

func (c *Client) StreamCount() (int64, error) {
	return c.r.Do("STREAMCOUNT").Int64()
}

func (c *Client) StreamCountEstimate() (int64, error) {
	return c.r.Do("STREAMCOUNTEST").Int64()
}
