package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"eventflowdb/graph/generated"
	"eventflowdb/graph/model"
	"fmt"
)

func (r *mutationResolver) Append(ctx context.Context, events []*model.EventData) ([]*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FromStream(ctx context.Context, stream *string, version *int, limit *int) ([]*model.Event, error) {
	return nil, nil
}

func (r *queryResolver) FromAllStreams(ctx context.Context, offset *string, limit *int) ([]*model.Event, error) {
	return nil, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) LoadFromStream(ctx context.Context, stream *string, version *int, limit *int) ([]*model.Event, error) {
	return nil, nil
}
