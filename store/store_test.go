package store

import (
	"math/rand"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/bbolt"
)

func GenerateRandomEvents(count int) []EventData {
	var events []EventData

	types := []string{
		"UserAdded", "ProductAdded", "OrderRefunded", "OrderCanceled", "ProductDiscontinued", "UserLocked",
	}

	for i := 0; i < count; i++ {
		t := types[rand.Intn(len(types))]

		payload := make([]byte, rand.Intn(128)+16)

		rand.Read(payload)

		events = append(events, EventData{
			Type: t,
			Data: payload,
		})
	}

	return events
}

func TestAppendingEvents(t *testing.T) {
	assert := assert.New(t)

	db, err := bbolt.Open("/tmp/events.db", 0666, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store, err := NewEventStore(db)
	if err != nil {
		t.Fatal(err)
	}

	stream := uuid.New()

	count := rand.Intn(9) + 1

	if _, err := store.AppendToStream(stream, 0, GenerateRandomEvents(count)); err != nil {
		t.Fatal(err)
	}

	totalEvents, err := store.GetTotalEvents()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(totalEvents, count)

	os.Remove("/tmp/events.db")
}
