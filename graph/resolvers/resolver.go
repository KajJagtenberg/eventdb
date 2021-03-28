package graph

import "github.com/hashicorp/raft"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	raft *raft.Raft
}
