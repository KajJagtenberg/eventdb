package model

import "time"

type Event struct {
	ID       string    `json:"id"`
	Stream   string    `json:"stream"`
	Version  int       `json:"version"`
	Type     string    `json:"type"`
	Data     string    `json:"data"`
	Metadata string    `json:"metadata"`
	AddedAt  time.Time `json:"added_at"`
}
