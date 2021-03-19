package store_test

import (
	"eventflowdb/store"
	"eventflowdb/util"
	"math/rand"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/bbolt"
)

func TestAppendingEvents(t *testing.T) {
	assert := assert.New(t)

	db, err := bbolt.Open("/tmp/events.db", 0666, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store, err := store.NewEventStore(db)
	if err != nil {
		t.Fatal(err)
	}

	stream := uuid.New()

	count := rand.Intn(9) + 1

	if _, err := store.AppendToStream(stream, 0, util.GenerateRandomEvents(count)); err != nil {
		t.Fatal(err)
	}

	totalEvents, err := store.GetTotalEvents()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(totalEvents, count)

	os.Remove("/tmp/events.db")
}
