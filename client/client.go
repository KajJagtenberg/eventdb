package client

import (
	"github.com/go-redis/redis"
	"github.com/oklog/ulid"
)

var (
	SINCE_START = ulid.ULID{}
)

type Config struct {
	Address string
}

type Client struct {
	r *redis.Client
}

func NewClient(config *Config) *Client {
	r := redis.NewClient(&redis.Options{
		Addr: config.Address,
	})

	return &Client{r}
}
