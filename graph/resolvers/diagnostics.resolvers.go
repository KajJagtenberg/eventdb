package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/kajjagtenberg/eventflowdb/graph/model"
)

var (
	start = time.Now()
)

func (r *queryResolver) Diagnostics(ctx context.Context) (*model.Diagnostics, error) {
	return &model.Diagnostics{
		Uptime:   int(time.Now().Sub(start).Seconds()),
		UptimeMs: int(time.Now().Sub(start).Milliseconds()),
	}, nil
}
