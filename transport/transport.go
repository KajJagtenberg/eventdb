package api

import (
	"context"
	"time"

	"github.com/hashicorp/raft"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/constants"
	"github.com/kajjagtenberg/eventflowdb/fsm"
	"google.golang.org/protobuf/proto"
)

type EventStoreService struct {
	raft *raft.Raft
}

func (s *EventStoreService) handle(op string, req proto.Message) (interface{}, error) {
	payload, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	cmd, err := proto.Marshal(&api.Command{
		Op:      op,
		Payload: payload,
	})
	if err != nil {
		return nil, err
	}

	future := s.raft.Apply(cmd, 500*time.Millisecond)
	if err := future.Error(); err != nil {
		return nil, err
	}

	response := future.Response().(fsm.ApplyResponse)

	return response.Data, response.Error
}

func (s *EventStoreService) Add(ctx context.Context, in *api.AddRequest) (*api.EventResponse, error) {
	res, err := s.handle("ADD", in)
	return res.(*api.EventResponse), err
}

func (s *EventStoreService) Get(ctx context.Context, in *api.GetRequest) (*api.EventResponse, error) {
	return nil, nil
}
func (s *EventStoreService) GetAll(ctx context.Context, in *api.GetAllRequest) (*api.EventResponse, error) {
	return nil, nil
}

func (s *EventStoreService) Checksum(context.Context, *api.ChecksumRequest) (*api.ChecksumResponse, error) {
	return nil, nil
}

func (s *EventStoreService) EventCount(context.Context, *api.EventCountRequest) (*api.EventCountResponse, error) {
	return nil, nil
}

func (s *EventStoreService) EventCountEstimate(context.Context, *api.EventCountRequest) (*api.EventCountResponse, error) {
	return nil, nil
}

func (s *EventStoreService) StreamCount(context.Context, *api.StreamCountRequest) (*api.StreamCountResponse, error) {
	return nil, nil
}

func (s *EventStoreService) StreamCountEstimate(context.Context, *api.StreamCountRequest) (*api.StreamCountResponse, error) {
	return nil, nil
}

func (s *EventStoreService) ListStreams(ctx context.Context, in *api.ListStreamsRequest) (*api.ListStreamsReponse, error) {
	return nil, nil
}

func (s *EventStoreService) Size(context.Context, *api.SizeRequest) (*api.SizeResponse, error) {
	return nil, nil
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

func NewEventStoreService(raft *raft.Raft) *EventStoreService {
	return &EventStoreService{raft: raft}
}
