package store

import (
	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

type Stream struct {
	Name   uuid.UUID   `json:"name"`
	Events []ulid.ULID `json:"events"`
}
