package transport

import (
	"context"
	"time"

	"github.com/eventflowdb/eventflowdb/api"
	"github.com/eventflowdb/eventflowdb/constants"
	"github.com/eventflowdb/eventflowdb/store"
)

type EventStore struct {
	api.UnimplementedEventStoreServer
	eventstore store.EventStore
}

func (s *EventStore) Add(ctx context.Context, req *api.AddRequest) (*api.EventResponse, error) {
	return s.eventstore.Add(req)
}

func (s *EventStore) Get(ctx context.Context, req *api.GetRequest) (*api.EventResponse, error) {
	return s.eventstore.Get(req)
}
func (s *EventStore) GetAll(ctx context.Context, req *api.GetAllRequest) (*api.EventResponse, error) {
	return s.eventstore.GetAll(req)
}

func (s *EventStore) EventCount(ctx context.Context, req *api.EventCountRequest) (*api.EventCountResponse, error) {
	return s.eventstore.EventCount(req)
}

func (s *EventStore) EventCountEstimate(ctx context.Context, req *api.EventCountEstimateRequest) (*api.EventCountResponse, error) {
	return s.eventstore.EventCountEstimate(req)
}

func (s *EventStore) StreamCount(ctx context.Context, req *api.StreamCountRequest) (*api.StreamCountResponse, error) {
	return s.eventstore.StreamCount(req)
}

func (s *EventStore) StreamCountEstimate(ctx context.Context, req *api.StreamCountEstimateRequest) (*api.StreamCountResponse, error) {
	return s.eventstore.StreamCountEstimate(req)
}

func (s *EventStore) ListStreams(ctx context.Context, req *api.ListStreamsRequest) (*api.ListStreamsReponse, error) {
	return s.eventstore.ListStreams(req)
}

func (s *EventStore) Size(ctx context.Context, req *api.SizeRequest) (*api.SizeResponse, error) {
	return s.eventstore.Size(req)
}

var (
	start = time.Now()
)

func (s *EventStore) Uptime(context.Context, *api.UptimeRequest) (*api.UptimeResponse, error) {
	uptime := time.Since(start)

	return &api.UptimeResponse{
		Uptime:      uptime.Milliseconds(),
		UptimeHuman: uptime.String(),
	}, nil
}

func (s *EventStore) Version(context.Context, *api.VersionRequest) (*api.VersionResponse, error) {
	return &api.VersionResponse{
		Version: constants.Version,
	}, nil
}

func NewEventStore(eventstore store.EventStore) *EventStore {
	return &EventStore{eventstore: eventstore}
}
