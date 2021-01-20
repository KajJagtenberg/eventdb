package store

import (
	"github.com/oklog/ulid"
)

type Stream struct {
	Version int
	Events  []ulid.ULID
}
