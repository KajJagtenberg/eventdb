package commands

import (
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/kajjagtenberg/go-commando"
)

const (
	CMD_EVENT_COUNT           = "eventcount"
	CMD_EVENT_COUNT_SHORT     = "ec"
	CMD_EVENT_COUNT_EST       = "eventcountest"
	CMD_EVENT_COUNT_EST_SHORT = "ece"
)

type EventCountResponse struct {
	Count int64 `json:"count"`
}

func EventCountHandler(store store.EventStore) commando.CommandHandler {
	return func(cmd commando.Command) (interface{}, error) {
		count, err := store.EventCount()
		if err != nil {
			return nil, err
		}

		return EventCountResponse{count}, nil
	}
}

func EventCountEstimateHandler(store store.EventStore) commando.CommandHandler {
	return func(cmd commando.Command) (interface{}, error) {
		count, err := store.EventCountEstimate()
		if err != nil {
			return nil, err
		}

		return EventCountResponse{count}, nil
	}
}

func SetupEventCounterHandler(dispatcher *commando.CommandDispatcher, store store.EventStore) {
	dispatcher.Register(CMD_EVENT_COUNT, CMD_EVENT_COUNT_SHORT, EventCountHandler(store))
	dispatcher.Register(CMD_EVENT_COUNT_EST, CMD_EVENT_COUNT_EST_SHORT, EventCountEstimateHandler(store))
}
