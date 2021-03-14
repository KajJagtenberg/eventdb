package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"eventflowdb/graph/generated"
	"eventflowdb/graph/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

func (r *mutationResolver) Append(ctx context.Context, stream string, version int, events []*model.EventData) ([]*model.RecordedEvent, error) {
	return nil, fmt.Errorf("not implemented")
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
	return nil, fmt.Errorf("not implemented")
}

func (r *queryResolver) LoadFromStream(ctx context.Context, stream string, skip int, limit int) ([]*model.RecordedEvent, error) {
	parsedStream, err := uuid.Parse(stream)
	if err != nil {
		return nil, err
	}

	records, err := r.EventStore.LoadFromStream(parsedStream, skip, limit)
	if err != nil {
		return nil, err
	}

	var result []*model.RecordedEvent

	for _, record := range records {
		result = append(result, &model.RecordedEvent{
			ID:       record.ID.String(),
			Stream:   record.Stream.String(),
			Version:  record.Version,
			Type:     record.Type,
			Data:     string(record.Data),
			Metadata: string(record.Metadata),
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
			Version:  record.Version,
			Type:     record.Type,
			Data:     string(record.Data),
			Metadata: string(record.Metadata),
			AddedAt:  record.AddedAt,
		})
	}

	return result, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
