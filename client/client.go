package client

import (
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/api"
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

func (c *Client) Add(stream uuid.UUID, version uint32, data []store.EventData) ([]store.Event, error) {
	req := api.AddRequest{
		Stream:  stream,
		Version: version,
		Events:  data,
	}

	cmd, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	response, err := c.r.Do("ADD", string(cmd)).Result()
	if err != nil {
		return nil, err
	}

	var events []store.Event

	for _, entry := range response.([]interface{}) {
		var event store.Event
		if err := json.Unmarshal([]byte(entry.(string)), &event); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (c *Client) Get(stream uuid.UUID, version uint32, limit uint32) ([]store.Event, error) {
	req := api.GetRequest{
		Stream:  stream,
		Version: version,
		Limit:   limit,
	}

	args, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := c.r.Do("GET", args).Result()
	if err != nil {
		return nil, err
	}

	var events []store.Event

	for _, entry := range response.([]interface{}) {
		var event store.Event
		if err := json.Unmarshal([]byte(entry.(string)), &event); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (c *Client) GetAll(offset ulid.ULID, limit uint32) ([]store.Event, error) {
	req := api.GetAllRequest{
		Offset: offset,
		Limit:  limit,
	}

	args, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := c.r.Do("GETALL", args).Result()
	if err != nil {
		return nil, err
	}

	var events []store.Event

	for _, entry := range response.([]interface{}) {
		var event store.Event
		if err := json.Unmarshal([]byte(entry.(string)), &event); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func NewClient(config *Config) *Client {
	r := redis.NewClient(&redis.Options{
		Addr: config.Address,
	})

	return &Client{r}
}
