package store

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID            string      `json:"id"`
	Stream        uuid.UUID   `json:"stream"`
	Version       int         `json:"version"`
	Type          string      `json:"type"`
	Data          interface{} `json:"data"`
	CausationID   string      `json:"causation_id"`
	CorrelationID string      `json:"correlation_id"`
	Timestamp     time.Time   `json:"ts"`
}

type AppendEvent struct {
	Type          string      `json:"type" validate:"required,ascii"`
	Data          interface{} `json:"data" validate:"required"`
	CausationID   string      `json:"causation_id" validate:"required,uuid"`
	CorrelationID string      `json:"correlation_id" validate:"required,uuid"`
}
