package transport

import (
	"context"
	"time"

	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/constants"
	"github.com/kajjagtenberg/eventflowdb/store"
)

type EventStoreService struct {
	api.UnimplementedEventStoreServiceServer
	eventstore store.EventStore
}

func (s *EventStoreService) Add(ctx context.Context, req *api.AddRequest) (*api.EventResponse, error) {
	return s.eventstore.Add(req)
}

func (s *EventStoreService) Get(ctx context.Context, req *api.GetRequest) (*api.EventResponse, error) {
	return s.eventstore.Get(req)
}
func (s *EventStoreService) GetAll(ctx context.Context, req *api.GetAllRequest) (*api.EventResponse, error) {
	return s.eventstore.GetAll(req)
}

func (s *EventStoreService) EventCount(ctx context.Context, req *api.EventCountRequest) (*api.EventCountResponse, error) {
	return s.eventstore.EventCount(req)
}

func (s *EventStoreService) EventCountEstimate(ctx context.Context, req *api.EventCountEstimateRequest) (*api.EventCountResponse, error) {
	return s.eventstore.EventCountEstimate(req)
}

func (s *EventStoreService) StreamCount(ctx context.Context, req *api.StreamCountRequest) (*api.StreamCountResponse, error) {
	return s.eventstore.StreamCount(req)
}

func (s *EventStoreService) StreamCountEstimate(ctx context.Context, req *api.StreamCountEstimateRequest) (*api.StreamCountResponse, error) {
	return s.eventstore.StreamCountEstimate(req)
}

func (s *EventStoreService) ListStreams(ctx context.Context, req *api.ListStreamsRequest) (*api.ListStreamsReponse, error) {
	return s.eventstore.ListStreams(req)
}

func (s *EventStoreService) Size(ctx context.Context, req *api.SizeRequest) (*api.SizeResponse, error) {
	return s.eventstore.Size(req)
}

var (
	start = time.Now()
)

func (s *EventStoreService) Uptime(context.Context, *api.UptimeRequest) (*api.UptimeResponse, error) {
	uptime := time.Since(start)

	return &api.UptimeResponse{
		Uptime:      uptime.Milliseconds(),
		UptimeHuman: uptime.String(),
	}, nil
}

func (s *EventStoreService) Version(context.Context, *api.VersionRequest) (*api.VersionResponse, error) {
	return &api.VersionResponse{
		Version: constants.Version,
	}, nil
}

func NewEventStoreService(eventstore store.EventStore) *EventStoreService {
	return &EventStoreService{eventstore: eventstore}
}
