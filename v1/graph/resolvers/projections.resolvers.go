package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/base64"
	"eventflowdb/graph/generated"
	"eventflowdb/graph/model"
	"fmt"
)

func (r *mutationResolver) CreateProjection(ctx context.Context, input *model.CreateProjection) (*model.Projection, error) {
	decoded, err := base64.StdEncoding.DecodeString(input.Code)
	if err != nil {
		return nil, err
	}

	proj, err := r.ProjectionEngine.CreateProjection(input.Name, string(decoded))
	if err != nil {
		return nil, err
	}

	return &model.Projection{
		ID:         proj.ID.String(),
		Name:       proj.Name,
		Code:       proj.Code,
		CreatedAt:  proj.CreatedAt,
		UpdatedAt:  proj.UpdatedAt,
		Checkpoint: proj.Checkpoint.String(),
	}, nil
}

func (r *queryResolver) Projections(ctx context.Context, skip int, limit int) ([]*model.Projection, error) {
	projections, err := r.ProjectionEngine.GetProjections()
	if err != nil {
		return nil, err
	}
	var result []*model.Projection

	for _, projection := range projections {
		result = append(result, &model.Projection{
			ID:         projection.ID.String(),
			Name:       projection.Name,
			Code:       projection.Code,
			CreatedAt:  projection.CreatedAt,
			UpdatedAt:  projection.UpdatedAt,
			Checkpoint: projection.Checkpoint.String(),
		})
	}

	return result, nil
}

func (r *queryResolver) Projection(ctx context.Context, id string) (*model.Projection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
