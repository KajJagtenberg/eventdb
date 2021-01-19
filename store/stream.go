package store

import (
	"github.com/oklog/ulid"
)

type Stream struct {
	Events []ulid.ULID `json:"events"`
}
