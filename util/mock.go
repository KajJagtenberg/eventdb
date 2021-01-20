package util

import (
	"eventdb/store"
	"log"
	"time"

	"github.com/google/uuid"
)

func AddMockEvents(eventstore *store.Store, count int) float64 {
	start := time.Now()

	for i := 0; i < count; i++ {
		cause := uuid.New().String()

		event := store.AppendEvent{
			Type:          "person_added",
			CausationID:   cause,
			CorrelationID: cause,
			Data: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{
				Name: "Kaj Jagtenberg",
				Age:  22,
			},
		}

		stream := uuid.New()

		if err := eventstore.AppendToStream(stream, 0, []store.AppendEvent{event}); err != nil {
			log.Fatal(err)
		}
	}

	passed := time.Now().Sub(start)

	return float64(count) / (passed.Seconds() + 1)
}
