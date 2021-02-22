// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Event struct {
	ID       string    `json:"id"`
	Stream   string    `json:"stream"`
	Version  int       `json:"version"`
	Type     string    `json:"type"`
	Data     string    `json:"data"`
	Metadata string    `json:"metadata"`
	AddedAt  time.Time `json:"added_at"`
}

type EventData struct {
	Type     string `json:"type"`
	Data     string `json:"data"`
	Metadata string `json:"metadata"`
}

type Stream struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}