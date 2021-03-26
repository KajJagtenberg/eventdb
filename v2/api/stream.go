package api

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	ErrNotImplemented = grpc.Errorf(codes.Unimplemented, "Not implemented")
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
