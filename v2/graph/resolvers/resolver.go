package resolvers

import (
	"time"

	"github.com/kajjagtenberg/eventflowdb/store"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Storage *store.Storage
	Start   time.Time
}
