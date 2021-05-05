package client

import (
	"errors"

	"github.com/go-redis/redis"
	"github.com/oklog/ulid"
)

var (
	SINCE_START = ulid.ULID{}
)

var (
	ErrInvalidResponse = errors.New("invalid response")
)

type Config struct {
	Address string
}

type Client struct {
	r *redis.Client
}

func (c *Client) Close() error {
	return c.r.Close()
}

func NewClient(config *Config) *Client {
	r := redis.NewClient(&redis.Options{
		Addr: config.Address,
	})

	return &Client{r}
}
