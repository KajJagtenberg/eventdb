package util

import (
	"eventdb/store"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func AddMockEvents(eventstore *store.Store, count int) float64 {
	start := time.Now()

	for count > 0 {
		events := []store.AppendEvent{}

		for i := 0; i < rand.Intn(9)+1; i++ {
			cause := uuid.New().String()

			events = append(events, store.AppendEvent{
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
			})

			count--
		}

		stream := uuid.New()

		if err := eventstore.AppendToStream(stream, 0, events); err != nil {
			log.Fatal(err)
		}
	}

	passed := time.Now().Sub(start)

	return float64(count) / (passed.Seconds() + 1)
}
