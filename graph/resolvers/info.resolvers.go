package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"eventflowdb/constants"
	"eventflowdb/graph/generated"
	"eventflowdb/graph/model"
	"time"
)

func (r *queryResolver) Info(ctx context.Context) (*model.Info, error) {
	return &model.Info{
		Name:    constants.Name,
		Version: constants.Version,
		Time:    int(time.Now().UTC().UnixNano()),
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
