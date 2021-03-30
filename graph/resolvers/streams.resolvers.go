package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/graph/model"
)

func (r *queryResolver) Streams(ctx context.Context, skip int, limit int) ([]string, error) {
	return r.persistence.Streams(skip, limit)
}

func (r *queryResolver) EventsFromStream(ctx context.Context, stream string, version int, limit int) ([]*model.Event, error) {
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

func (r *queryResolver) EventsFromLog(ctx context.Context, offset string, limit int) ([]*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}
