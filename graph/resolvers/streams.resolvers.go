package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
)

func (r *queryResolver) Streams(ctx context.Context, skip int, limit int) ([]string, error) {
	return r.persistence.Streams(skip, limit)
}
