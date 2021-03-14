package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"eventflowdb/graph/model"
	"fmt"
)

func (r *queryResolver) Projections(ctx context.Context, skip int, limit int) ([]*model.Projection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Projection(ctx context.Context, id string) (*model.Projection, error) {
	panic(fmt.Errorf("not implemented"))
}
