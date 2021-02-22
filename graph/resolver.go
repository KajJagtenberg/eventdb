package graph

import (
	"eventflowdb/store"

	"go.etcd.io/bbolt"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	EventStore *store.EventStore
	DB         *bbolt.DB
}
