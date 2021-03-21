package resolvers

import (
	"time"

	"github.com/kajjagtenberg/eventflowdb/cluster"
	"github.com/kajjagtenberg/eventflowdb/store"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Cluster *cluster.Cluster
	Storage *store.Storage
	Start   time.Time
}
