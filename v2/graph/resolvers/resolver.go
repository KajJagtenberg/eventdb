package resolvers

import (
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/kajjagtenberg/eventflowdb/store"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Memberlist *memberlist.Memberlist
	Storage    *store.Storage
	Start      time.Time
}
