package api

import (
	context "context"

	"github.com/kajjagtenberg/eventflowdb/store"
)

type StreamService struct {
	store store.Store
}

func (service *StreamService) AddEvents(context.Context, *AddEventsRequest) (*AddEventsResponse, error) {
	return nil, nil
}

func (service *StreamService) GetEvents(context.Context, *GetEventsRequest) (*GetEventsResponse, error) {
	return nil, nil
}

func (service *StreamService) LogEvents(context.Context, *LogEventsRequest) (*LogEventsResponse, error) {
	return nil, nil
}

func NewStreamService(store store.Store) *StreamService {
	return &StreamService{store}
}
