package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"eventflowdb/graph/model"
	"fmt"
)

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
