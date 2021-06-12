package transport

import (
	"context"
	"time"

	"github.com/hashicorp/raft"
	"github.com/kajjagtenberg/eventflowdb/constants"
)

type EventStoreService struct {
	UnimplementedEventStoreServiceServer
	raft *raft.Raft
}

func (s *EventStoreService) Add(ctx context.Context, in *AddRequest) (*EventResponse, error) {
	return nil, nil
}

func (s *EventStoreService) Get(ctx context.Context, in *GetRequest) (*EventResponse, error) {
	return nil, nil
}
func (s *EventStoreService) GetAll(ctx context.Context, in *GetAllRequest) (*EventResponse, error) {
	return nil, nil
}

func (s *EventStoreService) Checksum(context.Context, *ChecksumRequest) (*ChecksumResponse, error) {
	return nil, nil
}

func (s *EventStoreService) EventCount(context.Context, *EventCountRequest) (*EventCountResponse, error) {
	return nil, nil
}

func (s *EventStoreService) EventCountEstimate(context.Context, *EventCountRequest) (*EventCountResponse, error) {
	return nil, nil
}

func (s *EventStoreService) StreamCount(context.Context, *StreamCountRequest) (*StreamCountResponse, error) {
	return nil, nil
}

func (s *EventStoreService) StreamCountEstimate(context.Context, *StreamCountRequest) (*StreamCountResponse, error) {
	return nil, nil
}

func (s *EventStoreService) ListStreams(ctx context.Context, in *ListStreamsRequest) (*ListStreamsReponse, error) {
	return nil, nil
}

func (s *EventStoreService) Size(context.Context, *SizeRequest) (*SizeResponse, error) {
	return nil, nil
}

var (
	start = time.Now()
)

func (s *EventStoreService) Uptime(context.Context, *UptimeRequest) (*UptimeResponse, error) {
	uptime := time.Since(start)

	return &UptimeResponse{
		Uptime:      uptime.Milliseconds(),
		UptimeHuman: uptime.String(),
	}, nil
}

func (s *EventStoreService) Version(context.Context, *VersionRequest) (*VersionResponse, error) {
	return &VersionResponse{
		Version: constants.Version,
	}, nil
}

func NewEventStoreService(raft *raft.Raft) *EventStoreService {
	return &EventStoreService{raft: raft}
}
