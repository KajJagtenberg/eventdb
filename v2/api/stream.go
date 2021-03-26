package api

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	"github.com/kajjagtenberg/eventflowdb/cluster"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	ErrNotImplemented = grpc.Errorf(codes.Unimplemented, "Not implemented")
)

type StreamService struct {
	raft *raft.Raft
}

func (service *StreamService) AddEvents(ctx context.Context, req *AddEventsRequest) (*AddEventsResponse, error) {
	var stream uuid.UUID
	if err := stream.UnmarshalBinary(req.Stream); err != nil {
		return nil, err
	}

	version := req.Version

	if version < 0 {
		log.Println("Version is negative. This is a bug")
		return nil, errors.New("Version cannot be negative")
	}

	if len(req.Events) == 0 {
		return nil, errors.New("List of events cannot be empty")
	}

	var events []*cluster.AddCommand_Event

	for _, event := range req.Events {
		events = append(events, &cluster.AddCommand_Event{
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   event.CausationId,
			CorrelationId: event.CorrelationId,
		})
	}

	cmd, err := proto.Marshal(&cluster.ApplyLog{
		Command: &cluster.ApplyLog_Add{
			Add: &cluster.AddCommand{
				Stream:  stream[:],
				Version: version,
				Events:  events,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	future := service.raft.Apply(cmd, time.Second*5)
	if err := future.Error(); err != nil {
		return nil, err
	}

	// records := future.Response().([]*RecordedEvent)

	return &AddEventsResponse{
		// Events: records,
	}, nil
}

func (service *StreamService) GetEvents(ctx context.Context, req *GetEventsRequest) (*GetEventsResponse, error) {
	return nil, ErrNotImplemented
}

func (service *StreamService) LogEvents(ctx context.Context, req *LogEventsRequest) (*LogEventsResponse, error) {
	return nil, ErrNotImplemented
}

func NewStreamService(raft *raft.Raft) *StreamService {
	return &StreamService{raft}
}
