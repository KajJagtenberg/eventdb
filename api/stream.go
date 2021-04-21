package api

import (
	"bytes"
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
	stream, err := uuid.Parse(req.Stream)
	if err != nil {
		return nil, err
	}

	version := req.Version

	var events []store.EventData

	for _, event := range req.Events {
		var causationID ulid.ULID
		var correlationID ulid.ULID
		var err error

		if id := event.CausationId; len(id) != 0 {
			causationID, err = ulid.Parse(id)
			if err != nil {
				return nil, err
			}
		}

		if id := event.CorrelationId; len(id) != 0 {
			correlationID, err = ulid.Parse(id)
			if err != nil {
				return nil, err
			}
		}

		events = append(events, store.EventData{
			Type:          event.Type,
			Data:          bytes.NewBufferString(event.Data),
			Metadata:      bytes.NewBufferString(event.Metadata),
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
	stream, err := uuid.Parse(req.Stream)
	if err != nil {
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
	offset, err := ulid.Parse(req.Offset)
	if err != nil {
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

func (service *StreamService) StreamCount(ctx context.Context, req *StreamCountRequest) (*StreamCountResponse, error) {
	count, err := service.store.StreamCount()
	if err != nil {
		return nil, err
	}
	return &StreamCountResponse{Count: count}, nil
}

func (service *StreamService) EventCount(ctx context.Context, req *EventCountRequest) (*EventCountResponse, error) {
	count, err := service.store.EventCount()
	if err != nil {
		return nil, err
	}
	return &EventCountResponse{Count: count}, nil
}

/**/

func (service *StreamService) StreamCountEstimate(ctx context.Context, req *StreamCountEstimateRequest) (*StreamCountEstimateResponse, error) {
	count, err := service.store.StreamCountEstimate()
	if err != nil {
		return nil, err
	}
	return &StreamCountEstimateResponse{Count: count}, nil
}

func (service *StreamService) EventCountEstimate(ctx context.Context, req *EventCountEstimateRequest) (*EventCountEstimateResponse, error) {
	count, err := service.store.EventCountEstimate()
	if err != nil {
		return nil, err
	}
	return &EventCountEstimateResponse{Count: count}, nil
}

func NewStreamService(store store.Store) *StreamService {
	return &StreamService{store}
}

func mapEvents(in []store.Event) []*Event {
	var result []*Event
	for _, record := range in {
		result = append(result, &Event{
			Id:            record.ID.String(),
			Stream:        record.Stream.String(),
			Version:       record.Version,
			Type:          record.Type,
			Data:          record.Data.String(),
			Metadata:      record.Metadata.String(),
			CausationId:   record.CausationID.String(),
			CorrelationId: record.CorrelationID.String(),
			AddedAt:       record.AddedAt.UnixNano(),
		})
	}

	return result
}
