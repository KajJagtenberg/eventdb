package store

import (
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

type Stream struct {
	ID        uuid.UUID   `json:"name"`
	Events    []ulid.ULID `json:"events"`
	CreatedAt time.Time   `json:"created_at"`
}

func (stream *Stream) Size() int {
	return len(stream.Events)
}
