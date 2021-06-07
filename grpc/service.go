package grpc

import (
	context "context"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/oklog/ulid"
)

type EventServiceImpl struct {
	UnimplementedEventServiceServer
	store store.EventStore
}

func (s *EventServiceImpl) Add(ctx context.Context, in *AddRequest) (*EventResponse, error) {
	stream, err := uuid.Parse(in.Stream)
	if err != nil {
		return nil, err
	}

	var eventdata []store.EventData

	for _, e := range in.Events {
		causation_id, err := ulid.Parse(e.CausationId)
		if err != nil {
			return nil, err
		}
		correlation_id, err := ulid.Parse(e.CausationId)
		if err != nil {
			return nil, err
		}

		eventdata = append(eventdata, store.EventData{
			Type:          e.Type,
			Data:          e.Data,
			Metadata:      e.Metadata,
			CausationID:   causation_id,
			CorrelationID: correlation_id,
		})
	}

	events, err := s.store.Add(stream, in.Version, eventdata)
	if err != nil {
		return nil, err
	}

	var parsedEvent []*EventResponse_Event

	for _, event := range events {
		parsedEvent = append(parsedEvent, &EventResponse_Event{
			Id:            event.ID.String(),
			Stream:        event.Stream.String(),
			Version:       event.Version,
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   event.CausationID.String(),
			CorrelationId: event.CorrelationID.String(),
			AddedAt:       event.AddedAt.Unix(),
		})
	}

	return &EventResponse{
		Events: parsedEvent,
	}, nil
}
func (s *EventServiceImpl) Get(ctx context.Context, in *GetRequest) (*EventResponse, error) {
	stream, err := uuid.Parse(in.Stream)
	if err != nil {
		return nil, err
	}

	events, err := s.store.Get(stream, in.Version, in.Limit)
	if err != nil {
		return nil, err
	}

	var parsedEvent []*EventResponse_Event

	for _, event := range events {
		parsedEvent = append(parsedEvent, &EventResponse_Event{
			Id:            event.ID.String(),
			Stream:        event.Stream.String(),
			Version:       event.Version,
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   event.CausationID.String(),
			CorrelationId: event.CorrelationID.String(),
			AddedAt:       event.AddedAt.Unix(),
		})
	}

	return &EventResponse{
		Events: parsedEvent,
	}, nil
}
func (s *EventServiceImpl) GetAll(ctx context.Context, in *GetAllRequest) (*EventResponse, error) {
	offset, err := ulid.Parse(in.Offset)
	if err != nil {
		return nil, err
	}

	events, err := s.store.GetAll(offset, in.Limit)
	if err != nil {
		return nil, err
	}

	var parsedEvent []*EventResponse_Event

	for _, event := range events {
		parsedEvent = append(parsedEvent, &EventResponse_Event{
			Id:            event.ID.String(),
			Stream:        event.Stream.String(),
			Version:       event.Version,
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   event.CausationID.String(),
			CorrelationId: event.CorrelationID.String(),
			AddedAt:       event.AddedAt.Unix(),
		})
	}

	return &EventResponse{
		Events: parsedEvent,
	}, nil
}

func NewEventService(store store.EventStore) *EventServiceImpl {
	return &EventServiceImpl{store: store}
}
