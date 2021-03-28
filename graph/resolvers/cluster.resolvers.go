package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/kajjagtenberg/eventflowdb/graph/generated"
	"github.com/kajjagtenberg/eventflowdb/graph/model"
)

func (r *queryResolver) Cluster(ctx context.Context) (*model.Cluster, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
