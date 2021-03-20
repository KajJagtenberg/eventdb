package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
)

func (r *queryResolver) StreamCount(ctx context.Context) (int, error) {
	streamCount, err := r.Storage.StreamCount()
	if err != nil {
		return 0, err
	}

	return int(streamCount), nil
}
