package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"eventflowdb/constants"
	"eventflowdb/graph/generated"
	"time"
)

func (r *queryResolver) Version(ctx context.Context) (string, error) {
	return constants.Version, nil
}

func (r *queryResolver) Time(ctx context.Context) (int, error) {
	return time.Now().UTC().Nanosecond(), nil
}

func (r *queryResolver) Uptime(ctx context.Context) (int, error) {
	return int(time.Now().Sub(r.Startup).Seconds()), nil
}

func (r *queryResolver) TotalStreams(ctx context.Context) (int, error) {
	return r.EventStore.GetTotalStreams()
}

func (r *queryResolver) TotalEvents(ctx context.Context) (int, error) {
	return r.EventStore.GetTotalEvents()
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
