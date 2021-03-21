// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type ClusterNode struct {
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Address string `json:"address"`
}

type GetInput struct {
	Stream  string `json:"stream"`
	Version int    `json:"version"`
	Limit   int    `json:"limit"`
}

type LogInput struct {
	Offset string `json:"offset"`
	Limit  int    `json:"limit"`
}

type RecordedEvent struct {
	ID            string    `json:"id"`
	Stream        string    `json:"stream"`
	Version       int       `json:"version"`
	Type          string    `json:"type"`
	Data          string    `json:"data"`
	Metadata      string    `json:"metadata"`
	CausationID   string    `json:"causation_id"`
	CorrelationID string    `json:"correlation_id"`
	AddedAt       time.Time `json:"added_at"`
}

type Stream struct {
	ID      string    `json:"id"`
	AddedAt time.Time `json:"added_at"`
}

type StreamsInput struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}
