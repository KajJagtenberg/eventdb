package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/base64"
	"errors"
	"eventflowdb/graph/model"
	"eventflowdb/store"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

func (r *mutationResolver) Append(ctx context.Context, stream string, version int, events []*model.EventData) ([]*model.RecordedEvent, error) {
	streamId, err := uuid.Parse(stream)
	if err != nil {
		return nil, errors.New("Unable to parse stream ID")
	}

	var eventData []store.EventData

	for _, event := range events {
		data, err := base64.StdEncoding.DecodeString(event.Data)
		if err != nil {
			return nil, err
		}

		metadata, err := base64.StdEncoding.DecodeString(event.Metadata)
		if err != nil {
			return nil, err
		}

		eventData = append(eventData, store.EventData{
			Type:     event.Type,
			Data:     data,
			Metadata: metadata,
		})
	}

	records, err := r.EventStore.AppendToStream(streamId, version, eventData)
	if err != nil {
		return nil, err
	}

	var result []*model.RecordedEvent

	for _, record := range records {
		result = append(result, &model.RecordedEvent{
			ID:       record.ID.String(),
			Stream:   record.Stream.String(),
			Version:  int(record.Version),
			Type:     record.Type,
			Data:     string(record.Data),
			Metadata: string(record.Metadata),
			AddedAt:  record.AddedAt,
		})
	}

	return result, nil
}

func (r *queryResolver) Streams(ctx context.Context, skip int, limit int) ([]*model.Stream, error) {
	streams, err := r.EventStore.GetStreams(skip, limit)
	if err != nil {
		return nil, err
	}

	var result []*model.Stream

	for _, stream := range streams {
		result = append(result, &model.Stream{
			ID:        stream.ID.String(),
			Size:      stream.Size(),
			CreatedAt: stream.CreatedAt,
		})
	}

	return result, nil
}

func (r *queryResolver) Stream(ctx context.Context, id string) (*model.Stream, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	stream, err := r.EventStore.GetStream(parsedId)
	if err != nil {
		return nil, err
	}

	return &model.Stream{ID: stream.ID.String(), Size: stream.Size(), CreatedAt: stream.CreatedAt}, nil
}

func (r *queryResolver) LoadFromStream(ctx context.Context, stream string, skip int, limit int) ([]*model.RecordedEvent, error) {
	parsedId, err := uuid.Parse(stream)
	if err != nil {
		return nil, err
	}

	records, err := r.EventStore.LoadFromStream(parsedId, skip, limit)
	if err != nil {
		return nil, err
	}

	var result []*model.RecordedEvent

	for _, record := range records {
		result = append(result, &model.RecordedEvent{
			ID:       record.ID.String(),
			Stream:   record.Stream.String(),
			Version:  int(record.Version),
			Type:     record.Type,
			Data:     base64.StdEncoding.EncodeToString(record.Data),
			Metadata: base64.StdEncoding.EncodeToString(record.Metadata),
			AddedAt:  record.AddedAt,
		})
	}

	return result, nil
}

func (r *queryResolver) LoadFromAll(ctx context.Context, offset string, limit int) ([]*model.RecordedEvent, error) {
	var parsedOffset ulid.ULID
	var err error

	if len(offset) != 0 {
		parsedOffset, err = ulid.Parse(offset)
		if err != nil {
			return nil, err
		}
	}

	records, err := r.EventStore.LoadFromAll(parsedOffset, limit)
	if err != nil {
		return nil, err
	}

	var result []*model.RecordedEvent

	for _, record := range records {
		result = append(result, &model.RecordedEvent{
			ID:       record.ID.String(),
			Stream:   record.Stream.String(),
			Version:  int(record.Version),
			Type:     record.Type,
			Data:     base64.StdEncoding.EncodeToString(record.Data),
			Metadata: base64.StdEncoding.EncodeToString(record.Metadata),
			AddedAt:  record.AddedAt,
		})
	}

	return result, nil
}
