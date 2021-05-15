package commands

import "github.com/kajjagtenberg/eventflowdb/store"

const (
	CMD_EVENT_COUNT           = "eventcount"
	CMD_EVENT_COUNT_SHORT     = "ec"
	CMD_EVENT_COUNT_EST       = "eventcountest"
	CMD_EVENT_COUNT_EST_SHORT = "ece"
)

type EventCountResponse struct {
	Count int64 `json:"count"`
}

func EventCountHandler(store store.EventStore) CommandHandler {
	return func(cmd Command) (interface{}, error) {
		count, err := store.EventCount()
		if err != nil {
			return nil, err
		}

		return EventCountResponse{count}, nil
	}
}

func EventCountEstimateHandler(store store.EventStore) CommandHandler {
	return func(cmd Command) (interface{}, error) {
		count, err := store.EventCountEstimate()
		if err != nil {
			return nil, err
		}

		return EventCountResponse{count}, nil
	}
}
