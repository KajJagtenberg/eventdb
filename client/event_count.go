package client

func (c *Client) EventCount() (int64, error) {
	return c.r.Do("EVENTCOUNT").Int64()
}

func (c *Client) EventCountEstimate() (int64, error) {
	return c.r.Do("EVENTCOUNTEST").Int64()
}
