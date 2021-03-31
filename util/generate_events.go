package util

import (
	"crypto/rand"
	"log"
	"time"

	"github.com/kajjagtenberg/eventflowdb/api"
)

func GenerateEvents(n int) chan api.AddEventsRequest_EventData {
	result := make(chan api.AddEventsRequest_EventData)

	go func() {
		for i := 0; i < n; i++ {
			data := make([]byte, 100)
			metadata := make([]byte, 20)

			rand.Read(data)
			rand.Read(metadata)

			result <- api.AddEventsRequest_EventData{
				Type:     "Random",
				Data:     data,
				Metadata: metadata,
				AddedAt:  time.Now().UnixNano(),
			}
		}

		close(result)

		log.Println("Closed")
	}()

	return result
}
