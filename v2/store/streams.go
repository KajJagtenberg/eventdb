package store

import (
	"context"
)

type EventStoreService struct {
	storage *Storage
}

func (s *EventStoreService) Add(req *AddRequest, stream Streams_AddServer) error {
	events, err := s.storage.Add(req)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := stream.Send(event); err != nil {
			return err
		}
	}

	return nil
}

func (s *EventStoreService) Get(req *GetRequest, stream Streams_GetServer) error {
	events, err := s.storage.Get(req)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := stream.Send(event); err != nil {
			return err
		}
	}

	return nil
}

func (s *EventStoreService) Log(req *LogRequest, stream Streams_LogServer) error {
	events, err := s.storage.Log(req)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := stream.Send(event); err != nil {
			return err
		}
	}

	return nil
}

func (s *EventStoreService) StreamCount(ctx context.Context, req *StreamCountRequest) (*StreamCountResponse, error) {
	count, err := s.storage.StreamCount()
	if err != nil {
		return nil, err
	}

	return &StreamCountResponse{Count: uint64(count)}, nil
}

func (s *EventStoreService) EventCount(ctx context.Context, req *EventCountRequest) (*EventCountResponse, error) {
	count, err := s.storage.EventCount()
	if err != nil {
		return nil, err
	}

	return &EventCountResponse{Count: uint64(count)}, nil
}

func (s *EventStoreService) GetStreams(req *GetStreamsRequest, result Streams_GetStreamsServer) error {
	streams, err := s.storage.GetStreams(req.Skip, req.Limit)
	if err != nil {
		return err
	}

	for _, stream := range streams {
		if err := result.Send(stream); err != nil {
			return err
		}
	}

	return nil
}

func NewEventStoreService(storage *Storage) *EventStoreService {
	return &EventStoreService{storage}
}
