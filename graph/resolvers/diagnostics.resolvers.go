package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/kajjagtenberg/eventflowdb/graph/model"
)

func (r *queryResolver) Diagnostics(ctx context.Context) (*model.Diagnostics, error) {
	return &model.Diagnostics{
		Uptime:   int(time.Now().Sub(start).Seconds()),
		UptimeMs: int(time.Now().Sub(start).Milliseconds()),
	}, nil
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
var (
	start = time.Now()
)
