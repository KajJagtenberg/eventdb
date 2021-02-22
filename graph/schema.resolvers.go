package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/base64"
	"eventflowdb/graph/generated"
	"eventflowdb/graph/model"
	"fmt"

	"github.com/google/uuid"
)

func (r *eventResolver) Stream(ctx context.Context, obj *model.Event) (*model.Stream, error) {
	return nil, nil
}

func (r *mutationResolver) Append(ctx context.Context, stream string, version int, events []*model.EventData) (*model.Stream, error) {
	return nil, nil
}

func (r *queryResolver) Streams(ctx context.Context, skip *int, limit *int) (*model.Streams, error) {
	return nil, nil
}

func (r *queryResolver) Stream(ctx context.Context, id string) (*model.Stream, error) {
	return nil, nil
}

func (r *queryResolver) All(ctx context.Context, offset string, limit *int) ([]*model.Event, error) {
	return nil, nil
}

func (r *streamResolver) Events(ctx context.Context, obj *model.Stream, skip *int, limit *int) ([]*model.Event, error) {
	stream, err := uuid.Parse(obj.Name)
	if err != nil {
		return nil, err
	}

	records, _, err := r.EventStore.LoadFromStream(stream, *skip, *limit)
	if err != nil {
		return nil, err
	}

	var events []*model.Event

	for _, record := range records {
		events = append(events, &model.Event{
			ID:       record.ID.String(),
			Stream:   record.Stream.String(),
			Version:  record.Version,
			Type:     record.Type,
			Data:     base64.StdEncoding.EncodeToString(record.Data),
			Metadata: base64.StdEncoding.EncodeToString(record.Metadata),
			AddedAt:  record.AddedAt,
		})
	}

	return events, nil
}

func (r *streamsResolver) Streams(ctx context.Context, obj *model.Streams) ([]*model.Stream, error) {
	panic(fmt.Errorf("not implemented"))
}

// Event returns generated.EventResolver implementation.
func (r *Resolver) Event() generated.EventResolver { return &eventResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Stream returns generated.StreamResolver implementation.
func (r *Resolver) Stream() generated.StreamResolver { return &streamResolver{r} }

// Streams returns generated.StreamsResolver implementation.
func (r *Resolver) Streams() generated.StreamsResolver { return &streamsResolver{r} }

type eventResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type streamResolver struct{ *Resolver }
type streamsResolver struct{ *Resolver }
