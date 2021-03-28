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
	"github.com/kajjagtenberg/eventflowdb/persistence"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	ErrNotImplemented = grpc.Errorf(codes.Unimplemented, "Not implemented")
)

type StreamService struct {
	raft        *raft.Raft
	persistence *persistence.Persistence
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

	var events []*cluster.AddCommand_EventData

	for _, event := range req.Events {
		events = append(events, &cluster.AddCommand_EventData{
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   event.CausationId,
			CorrelationId: event.CorrelationId,
			AddedAt:       event.AddedAt,
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

	result := future.Response().(cluster.ApplyResult)

	if result.Error != nil {
		return nil, result.Error
	}

	var records []*Event

	for _, event := range result.Value.([]persistence.Event) {
		records = append(records, &Event{
			Id:            event.ID[:],
			Stream:        event.Stream[:],
			Version:       event.Version,
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   event.CausationID[:],
			CorrelationId: event.CorrelationID[:],
			AddedAt:       event.AddedAt.UnixNano(),
		})
	}

	return &AddEventsResponse{
		Events: records,
	}, nil
}

func (service *StreamService) GetEvents(ctx context.Context, req *GetEventsRequest) (*GetEventsResponse, error) {
	var stream uuid.UUID
	if err := stream.UnmarshalBinary(req.Stream); err != nil {
		return nil, err
	}

	version := req.Version
	limit := req.Limit

	events, err := service.persistence.Get(stream, version, limit)
	if err != nil {
		return nil, err
	}

	var records []*Event

	for _, event := range events {
		records = append(records, &Event{
			Id:            event.ID[:],
			Stream:        event.Stream[:],
			Version:       event.Version,
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   event.CausationID[:],
			CorrelationId: event.CorrelationID[:],
			AddedAt:       event.AddedAt.UnixNano(),
		})
	}

	return &GetEventsResponse{
		Events: records,
	}, nil
}

func (service *StreamService) LogEvents(ctx context.Context, req *LogEventsRequest) (*LogEventsResponse, error) {
	return nil, ErrNotImplemented
}

func NewStreamService(raft *raft.Raft, persistence *persistence.Persistence) *StreamService {
	return &StreamService{raft, persistence}
}
