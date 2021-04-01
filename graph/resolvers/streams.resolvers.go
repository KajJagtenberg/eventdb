package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/graph/model"
	"github.com/oklog/ulid"
)

func (r *queryResolver) Streams(ctx context.Context, skip int, limit int) ([]string, error) {
	return r.persistence.Streams(skip, limit)
}

func (r *queryResolver) Get(ctx context.Context, stream string, version int, limit int) ([]*model.Event, error) {
	streamId, err := uuid.Parse(stream)
	if err != nil {
		return nil, err
	}

	events, err := r.persistence.Get(streamId, uint32(version), uint32(limit))
	if err != nil {
		return nil, err
	}

	var result []*model.Event

	for _, event := range events {
		result = append(result, &model.Event{
			ID:            event.ID.String(),
			Stream:        event.Stream.String(),
			Version:       int(event.Version),
			Type:          event.Type,
			Data:          string(event.Data),
			Metadata:      string(event.Metadata),
			CausationID:   event.CausationID.String(),
			CorrelationID: event.CorrelationID.String(),
			AddedAt:       event.AddedAt,
		})
	}

	return result, nil
}

func (r *queryResolver) Log(ctx context.Context, offset string, limit int) ([]*model.Event, error) {
	offsetId, err := ulid.Parse(offset)
	if err != nil {
		return nil, err
	}

	events, err := r.persistence.Log(offsetId, uint32(limit))
	if err != nil {
		return nil, err
	}

	var result []*model.Event

	for _, event := range events {
		result = append(result, &model.Event{
			ID:            event.ID.String(),
			Stream:        event.Stream.String(),
			Version:       int(event.Version),
			Type:          event.Type,
			Data:          string(event.Data),
			Metadata:      string(event.Metadata),
			CausationID:   event.CausationID.String(),
			CorrelationID: event.CorrelationID.String(),
			AddedAt:       event.AddedAt,
		})
	}

	return result, nil
}

func (r *queryResolver) EventCount(ctx context.Context) (int, error) {
	count, err := r.persistence.EventCount()
	if err != nil {
		return 0, err
	}

	return int(count), err
}

func (r *queryResolver) Checksum(ctx context.Context) (string, error) {
	checksum, err := r.persistence.Checksum()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(checksum), nil
}
