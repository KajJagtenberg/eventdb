package client

import (
	"encoding/json"
	"log"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
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

func (c *Client) Log(offset ulid.ULID, limit uint32) ([]store.Event, error) {
	res, err := c.r.Do("log", offset.String(), limit).Result()
	if err != nil {
		return nil, err
	}

	log.Println(res)

	return nil, nil
}

func (c *Client) Add(stream uuid.UUID, version uint32, events []store.EventData) ([]store.Event, error) {
	serializedEvents, err := json.Marshal(events)
	if err != nil {
		return nil, err
	}

	response, err := c.r.Do("ADD", stream.String(), version, serializedEvents).String()
	if err != nil {
		return nil, err
	}

	var result []store.Event
	if err := json.Unmarshal([]byte(response), &events); err != nil {
		return nil, err
	}

	return result, nil
}

func NewClient(config *Config) *Client {
	r := redis.NewClient(&redis.Options{
		Addr: config.Address,
	})

	return &Client{r}
}
