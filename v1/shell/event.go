package shell

import (
	"time"
)

type Event struct {
	ID       string    `json:"id"`
	Stream   string    `json:"stream"`
	Version  int       `json:"version"`
	Type     string    `json:"type"`
	Data     []byte    `json:"data"`
	Metadata []byte    `json:"metadata"`
	AddedAt  time.Time `json:"added_at"`
}
