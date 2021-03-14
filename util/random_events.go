package util

import (
	"eventflowdb/store"
	"math/rand"
)

func GenerateRandomEvents(count int) []store.EventData {
	var events []store.EventData

	types := []string{
		"UserAdded", "ProductAdded", "OrderRefunded", "OrderCanceled", "ProductDiscontinued", "UserLocked",
	}

	for i := 0; i < count; i++ {
		t := types[rand.Intn(len(types))]

		payload := make([]byte, rand.Intn(128)+16)

		rand.Read(payload)

		events = append(events, store.EventData{
			Type: t,
			Data: payload,
		})
	}

	return events
}
