package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/kajjagtenberg/eventflowdb/graph/generated"
	"github.com/kajjagtenberg/eventflowdb/graph/model"
)

func (r *queryResolver) Healthscore(ctx context.Context) (int, error) {
	return r.Cluster.GetList().GetHealthScore(), nil
}

func (r *queryResolver) Nodes(ctx context.Context) ([]*model.ClusterNode, error) {
	nodes := []*model.ClusterNode{}

	for _, member := range r.Cluster.GetList().Members() {
		nodes = append(nodes, &model.ClusterNode{
			IP:      member.Addr.String(),
			Port:    int(member.Port),
			Address: member.Address(),
		})
	}

	return nodes, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
