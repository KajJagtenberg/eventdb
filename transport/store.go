package transport

import (
	"context"
	"encoding/base32"
	"time"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/oklog/ulid"
)

type EventStoreService struct {
	UnimplementedEventStoreServiceServer
	store store.EventStore
}

func (s *EventStoreService) Add(ctx context.Context, in *AddRequest) (*EventResponse, error) {
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
func (s *EventStoreService) Get(ctx context.Context, in *GetRequest) (*EventResponse, error) {
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
func (s *EventStoreService) GetAll(ctx context.Context, in *GetAllRequest) (*EventResponse, error) {
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

func (s *EventStoreService) Checksum(context.Context, *ChecksumRequest) (*ChecksumResponse, error) {
	id, sum, err := s.store.Checksum()
	if err != nil {
		return nil, err
	}

	return &ChecksumResponse{
		Id:       id.String(),
		Checksum: base32.StdEncoding.EncodeToString(sum),
	}, nil
}

func (s *EventStoreService) EventCount(context.Context, *EventCountRequest) (*EventCountResponse, error) {
	count, err := s.store.EventCount()
	if err != nil {
		return nil, err
	}

	return &EventCountResponse{
		Count: count,
	}, nil
}

func (s *EventStoreService) EventCountEstimate(context.Context, *EventCountRequest) (*EventCountResponse, error) {
	count, err := s.store.EventCountEstimate()
	if err != nil {
		return nil, err
	}

	return &EventCountResponse{
		Count: count,
	}, nil
}

func (s *EventStoreService) ListStreams(ctx context.Context, in *ListStreamsRequest) (*ListStreamsReponse, error) {
	streams, err := s.store.ListStreams(in.Skip, in.Limit)
	if err != nil {
		return nil, err
	}

	var result []*ListStreamsReponse_Stream

	for _, stream := range streams {
		var events []string

		for _, event := range stream.Events {
			events = append(events, event.String())
		}

		result = append(result, &ListStreamsReponse_Stream{
			Id:      stream.ID.String(),
			Events:  events,
			AddedAt: stream.AddedAt.Unix(),
		})
	}

	return &ListStreamsReponse{
		Streams: result,
	}, nil
}

func (s *EventStoreService) Size(context.Context, *SizeRequest) (*SizeResponse, error) {
	size, err := s.store.Size()
	if err != nil {
		return nil, err
	}

	return &SizeResponse{
		Size: size,
	}, nil
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

func NewEventStoreService(store store.EventStore) *EventStoreService {
	return &EventStoreService{store: store}
}
