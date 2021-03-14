package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"eventflowdb/graph/generated"
	"eventflowdb/graph/model"
	"fmt"
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
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) LoadFromAll(ctx context.Context, offset string, limit int) ([]*model.RecordedEvent, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
