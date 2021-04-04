package api

import (
	context "context"
	"time"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/oklog/ulid"
)

type StreamService struct {
	store store.Store
}

func (service *StreamService) AddEvents(ctx context.Context, req *AddEventsRequest) (*AddEventsResponse, error) {
	var stream uuid.UUID
	if err := stream.UnmarshalBinary(req.Stream); err != nil {
		return nil, err
	}

	version := req.Version

	var events []store.EventData

	for _, event := range req.Events {
		var causationID ulid.ULID
		var correlationID ulid.ULID

		if id := event.CausationId; id != nil {
			if err := causationID.UnmarshalBinary(id); err != nil {
				return nil, err
			}
		}

		if id := event.CorrelationId; id != nil {
			if err := correlationID.UnmarshalBinary(id); err != nil {
				return nil, err
			}
		}

		events = append(events, store.EventData{
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationID:   causationID,
			CorrelationID: correlationID,
			AddedAt:       time.Unix(0, event.AddedAt),
		})
	}

	records, err := service.store.Add(stream, version, events)
	if err != nil {
		return nil, err
	}

	result := &AddEventsResponse{
		Events: mapEvents(records),
	}

	return result, nil
}

func (service *StreamService) GetEvents(ctx context.Context, req *GetEventsRequest) (*GetEventsResponse, error) {
	var stream uuid.UUID
	if err := stream.UnmarshalBinary(req.Stream); err != nil {
		return nil, err
	}

	version := req.Version
	limit := req.Limit

	records, err := service.store.Get(stream, version, limit)
	if err != nil {
		return nil, err
	}

	result := &GetEventsResponse{
		Events: mapEvents(records),
	}

	return result, nil
}

func (service *StreamService) LogEvents(ctx context.Context, req *LogEventsRequest) (*LogEventsResponse, error) {
	var offset ulid.ULID
	if err := offset.UnmarshalBinary(req.Offset); err != nil {
		return nil, err
	}

	limit := req.Limit

	records, err := service.store.Log(offset, limit)
	if err != nil {
		return nil, err
	}

	result := &LogEventsResponse{
		Events: mapEvents(records),
	}

	return result, nil
}

func NewStreamService(store store.Store) *StreamService {
	return &StreamService{store}
}

func mapEvents(in []store.Event) []*Event {
	var result []*Event
	for _, record := range in {
		result = append(result, &Event{
			Id:            record.ID[:],
			Stream:        record.Stream[:],
			Version:       record.Version,
			Type:          record.Type,
			Data:          record.Data,
			Metadata:      record.Metadata,
			CausationId:   record.CausationID[:],
			CorrelationId: record.CorrelationID[:],
			AddedAt:       record.AddedAt.UnixNano(),
		})
	}

	return result
}
