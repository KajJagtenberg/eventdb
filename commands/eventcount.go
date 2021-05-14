package commands

import "github.com/kajjagtenberg/eventflowdb/store"

type EventCountResponse struct {
	Count int64 `json:"count"`
}

func EventCountHandler(store store.Store) CommandHandler {
	return func(cmd Command) (interface{}, error) {
		count, err := store.EventCount()
		if err != nil {
			return nil, err
		}

		return EventCountResponse{count}, nil
	}
}

func EventCountEstimateHandler(store store.Store) CommandHandler {
	return func(cmd Command) (interface{}, error) {
		count, err := store.EventCountEstimate()
		if err != nil {
			return nil, err
		}

		return EventCountResponse{count}, nil
	}
}
