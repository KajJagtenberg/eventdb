package transport

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
	api.UnimplementedEventStoreServiceServer
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

	response := future.Response().(*fsm.ApplyResponse)

	return response.Data, response.Error
}

func (s *EventStoreService) Add(ctx context.Context, req *api.AddRequest) (*api.EventResponse, error) {
	res, err := s.handle("ADD", req)
	if err != nil {
		return nil, err
	}
	return res.(*api.EventResponse), err
}

func (s *EventStoreService) Get(ctx context.Context, req *api.GetRequest) (*api.EventResponse, error) {
	res, err := s.handle("GET", req)
	if err != nil {
		return nil, err
	}
	return res.(*api.EventResponse), err
}
func (s *EventStoreService) GetAll(ctx context.Context, req *api.GetAllRequest) (*api.EventResponse, error) {
	res, err := s.handle("GET_ALL", req)
	if err != nil {
		return nil, err
	}
	return res.(*api.EventResponse), err
}

func (s *EventStoreService) EventCount(ctx context.Context, req *api.EventCountRequest) (*api.EventCountResponse, error) {
	res, err := s.handle("EVENT_COUNT", req)
	if err != nil {
		return nil, err
	}
	return res.(*api.EventCountResponse), err
}

func (s *EventStoreService) EventCountEstimate(ctx context.Context, req *api.EventCountEstimateRequest) (*api.EventCountResponse, error) {
	res, err := s.handle("EVENT_COUNT_ESTIMATE", req)
	if err != nil {
		return nil, err
	}
	return res.(*api.EventCountResponse), err
}

func (s *EventStoreService) StreamCount(ctx context.Context, req *api.StreamCountRequest) (*api.StreamCountResponse, error) {
	res, err := s.handle("STREAM_COUNT", req)
	if err != nil {
		return nil, err
	}
	return res.(*api.StreamCountResponse), err
}

func (s *EventStoreService) StreamCountEstimate(ctx context.Context, req *api.StreamCountEstimateRequest) (*api.StreamCountResponse, error) {
	res, err := s.handle("STREAM_COUNT_ESTIMATE", req)
	if err != nil {
		return nil, err
	}
	return res.(*api.StreamCountResponse), err
}

func (s *EventStoreService) ListStreams(ctx context.Context, req *api.ListStreamsRequest) (*api.ListStreamsReponse, error) {
	res, err := s.handle("LIST_STREAMS", req)
	if err != nil {
		return nil, err
	}
	return res.(*api.ListStreamsReponse), err
}

func (s *EventStoreService) Size(ctx context.Context, req *api.SizeRequest) (*api.SizeResponse, error) {
	res, err := s.handle("SIZE", req)
	if err != nil {
		return nil, err
	}
	return res.(*api.SizeResponse), nil
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
