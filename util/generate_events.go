package util

import (
	"crypto/rand"
	"time"

	"github.com/kajjagtenberg/eventflowdb/persistence"
)

func GenerateEvents(n int) chan persistence.EventData {
	result := make(chan persistence.EventData)

	go func() {
		for i := 0; i < n; i++ {
			data := make([]byte, 100)
			metadata := make([]byte, 20)

			rand.Read(data)
			rand.Read(metadata)

			result <- persistence.EventData{
				Type:     "Random",
				Data:     data,
				Metadata: metadata,
				AddedAt:  time.Now(),
			}
		}
	}()

	return result
}
