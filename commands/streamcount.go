package commands

import (
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/kajjagtenberg/go-commando"
)

const (
	CMD_STREAM_COUNT           = "streamcount"
	CMD_STREAM_COUNT_SHORT     = "sc"
	CMD_STREAM_COUNT_EST       = "streamcountest"
	CMD_STREAM_COUNT_EST_SHORT = "sce"
)

type StreamCountResponse struct {
	Count int64 `json:"count"`
}

func StreamCountHandler(store store.EventStore) commando.CommandHandler {
	return func(cmd commando.Command) (interface{}, error) {
		count, err := store.StreamCount()
		if err != nil {
			return nil, err
		}

		return StreamCountResponse{count}, nil
	}
}

func StreamCountEstimateHandler(store store.EventStore) commando.CommandHandler {
	return func(cmd commando.Command) (interface{}, error) {
		count, err := store.StreamCountEstimate()
		if err != nil {
			return nil, err
		}

		return StreamCountResponse{count}, nil
	}
}
