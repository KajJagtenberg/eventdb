package resolvers

import (
	"github.com/hashicorp/raft"
	"github.com/kajjagtenberg/eventflowdb/persistence"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	raft        *raft.Raft
	persistence *persistence.Persistence
}

func NewResolver(raft *raft.Raft, persistence *persistence.Persistence) *Resolver {
	return &Resolver{raft, persistence}
}
