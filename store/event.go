package store

import (
	"encoding/json"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

type Event struct {
	ID            ulid.ULID       `json:"id"`
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
	Timestamp     time.Time   `json:"ts"`
	CausationID   string      `json:"causation_id" validate:"required,uuid"`
	CorrelationID string      `json:"correlation_id" validate:"required,uuid"`
}

func (e *AppendEvent) Validate() error {
	validator := validator.New()

	// TODO: Add friendlier error messages

	return validator.Struct(e)
}
