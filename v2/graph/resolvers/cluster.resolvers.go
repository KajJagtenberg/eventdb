package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/kajjagtenberg/eventflowdb/graph/generated"
	"github.com/kajjagtenberg/eventflowdb/graph/model"
)

func (r *queryResolver) Cluster(ctx context.Context) (*model.Cluster, error) {
	list := r.Memberlist

	result := model.Cluster{
		Healthscore: list.GetHealthScore(),
	}

	for _, member := range list.Members() {
		result.Nodes = append(result.Nodes, &model.ClusterNode{
			IP:      member.Addr.String(),
			Address: member.Address(),
			Port:    int(member.Port),
		})
	}

	return &result, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
