package api

import (
	"context"
	"errors"
)

var (
	ErrNotImplemented = errors.New("Not implemented")
)

type StreamService struct{}

func (service *StreamService) AddEvents(context.Context, *AddEventsRequest) (*AddEventsResponse, error) {
	return nil, ErrNotImplemented
}
func (service *StreamService) GetEvents(context.Context, *GetEventsRequest) (*GetEventsResponse, error) {
	return nil, ErrNotImplemented
}
func (service *StreamService) LogEvents(context.Context, *LogEventsRequest) (*LogEventsResponse, error) {
	return nil, ErrNotImplemented
}

func NewStreamService() *StreamService {
	return &StreamService{}
}
