package store

import "time"

type Event struct {
	ID            string      `json:"id"`
	Stream        string      `json:"stream"`
	Version       int         `json:"version"`
	Type          string      `json:"type"`
	Data          interface{} `json:"data"`
	CausationID   string      `json:"causation_id"`
	CorrelationID string      `json:"correlation_id"`
	Timestamp     time.Time   `json:"ts"`
}

type AppendEvent struct {
	Type          string      `json:"type"`
	Data          interface{} `json:"data"`
	CausationID   string      `json:"causation_id"`
	CorrelationID string      `json:"correlation_id"`
}
