package store

import (
	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

type EventData struct {
	Type     string `json:"type"`
	Data     []byte `json:"data"`
	Metadata []byte `json:"metadata"`
}

type RecordedEvent struct {
	ID       ulid.ULID `json:"id"`
	Stream   uuid.UUID `json:"stream"`
	Version  int       `json:"version"`
	Type     string    `json:"type"`
	Data     []byte    `json:"data"`
	Metadata []byte    `json:"metadata"`
}
