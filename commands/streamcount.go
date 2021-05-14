package commands

import "github.com/kajjagtenberg/eventflowdb/store"

type StreamCountResponse struct {
	Count int64 `json:"count"`
}

func StreamCountHandler(store store.Store) CommandHandler {
	return func(cmd Command) (interface{}, error) {
		count, err := store.StreamCount()
		if err != nil {
			return nil, err
		}

		return StreamCountResponse{count}, nil
	}
}

func StreamCountEstimateHandler(store store.Store) CommandHandler {
	return func(cmd Command) (interface{}, error) {
		count, err := store.StreamCountEstimate()
		if err != nil {
			return nil, err
		}

		return StreamCountResponse{count}, nil
	}
}
