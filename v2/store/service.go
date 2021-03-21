package store

import "context"

type StoreService struct {
	storage *Storage
}

func (s *StoreService) Add(req *AddRequest, stream EventStore_AddServer) error {
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

func (s *StoreService) Get(req *GetRequest, stream EventStore_GetServer) error {
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

func (s *StoreService) Log(req *LogRequest, stream EventStore_LogServer) error {
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

func (s *StoreService) StreamCount(ctx context.Context, req *StreamCountRequest) (*StreamCountResponse, error) {
	count, err := s.storage.StreamCount()
	if err != nil {
		return nil, err
	}

	return &StreamCountResponse{Count: uint64(count)}, nil
}

func (s *StoreService) EventCount(ctx context.Context, req *EventCountRequest) (*EventCountResponse, error) {
	count, err := s.storage.EventCount()
	if err != nil {
		return nil, err
	}

	return &EventCountResponse{Count: uint64(count)}, nil
}

func (s *StoreService) GetStreams(req *GetStreamsRequest, result EventStore_GetStreamsServer) error {
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

func NewStoreService(storage *Storage) *StoreService {
	return &StoreService{storage}
}
