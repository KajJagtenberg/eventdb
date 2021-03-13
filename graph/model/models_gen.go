// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type EventData struct {
	Type     string `json:"type"`
	Data     string `json:"data"`
	Metadata string `json:"metadata"`
}

type RecordedEvent struct {
	ID       string    `json:"id"`
	Stream   string    `json:"stream"`
	Version  int       `json:"version"`
	Type     string    `json:"type"`
	Data     string    `json:"data"`
	Metadata string    `json:"metadata"`
	AddedAt  time.Time `json:"added_at"`
}

type Stream struct {
	Name      string           `json:"name"`
	Size      int              `json:"size"`
	Events    []*RecordedEvent `json:"events"`
	CreatedAt time.Time        `json:"created_at"`
}
