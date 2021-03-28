package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/kajjagtenberg/eventflowdb/graph/generated"
	"github.com/kajjagtenberg/eventflowdb/graph/model"
)

func (r *queryResolver) Cluster(ctx context.Context) (*model.Cluster, error) {

	configFuture := r.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		return nil, err
	}

	config := configFuture.Configuration()

	return &model.Cluster{
		Leader: string(r.raft.Leader()),
		Size:   len(config.Servers),
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
