package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"
)

var (
	start = time.Now()
)

func (r *queryResolver) Uptime(ctx context.Context) (int, error) {
	return int(time.Now().Sub(start).Seconds()), nil
}

func (r *queryResolver) UptimeMs(ctx context.Context) (int, error) {
	return int(time.Now().Sub(start).Milliseconds()), nil
}
