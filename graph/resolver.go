package graph

import "eventflowdb/store"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	EventStore *store.EventStore
}
