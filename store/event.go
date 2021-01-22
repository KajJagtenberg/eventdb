package store

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID            uuid.UUID       `json:"id"`
	Stream        uuid.UUID       `json:"stream"`
	Version       int             `json:"version"`
	Type          string          `json:"type"`
	Data          json.RawMessage `json:"data"`
	Metadata      interface{}     `json:"metadata"`
	CausationID   string          `json:"causation_id"`
	CorrelationID string          `json:"correlation_id"`
	Timestamp     time.Time       `json:"ts"`
}

type AppendEvent struct {
	Type          string      `json:"type" validate:"required,ascii"`
	Data          interface{} `json:"data" validate:"required"`
	Metadata      interface{} `json:"metadata"`
	CausationID   string      `json:"causation_id" validate:"required,uuid"`
	CorrelationID string      `json:"correlation_id" validate:"required,uuid"`
}
