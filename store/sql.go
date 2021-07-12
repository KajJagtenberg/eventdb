package store

import (
	"github.com/eventflowdb/eventflowdb/api"
)

type SQLStore struct {
}

func (s *SQLStore) GetStream(*api.GetStreamRequest) (*api.GetStreamResponse, error) {
	return nil, nil
}

func (s *SQLStore) GetGlobalStream(*api.GetGlobalStreamRequest) (*api.GetGlobalStreamResponse, error) {
	return nil, nil
}

func (s *SQLStore) AppendToStream(*api.AppendToStreamRequest) (*api.AppendToStreamResponse, error) {
	return nil, nil
}

func (s *SQLStore) GetEvent(*api.GetEventRequest) (*api.Event, error) {
	return nil, nil
}

func (s *SQLStore) Size(*api.SizeRequest) (*api.SizeResponse, error) {
	return nil, nil
}

func (s *SQLStore) EventCount(*api.EventCountRequest) (*api.EventCountResponse, error) {
	return nil, nil
}

func (s *SQLStore) StreamCount(*api.StreamCountRequest) (*api.StreamCountResponse, error) {
	return nil, nil
}

func (s *SQLStore) ListStreams(*api.ListStreamsRequest) (*api.ListStreamsReponse, error) {
	return nil, nil
}

func (s *SQLStore) Close() error {
	return nil
}

func NewSQLStore() (*SQLStore, error) {
	return &SQLStore{}, nil
}
