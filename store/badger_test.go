package store_test

import (
	"crypto/rand"
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/oklog/ulid"
	"github.com/stretchr/testify/assert"
)

func BenchmarkAddMemory(t *testing.B) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	eventstore, err := store.NewBadgerEventStore(store.BadgerStoreOptions{
		DB:             db,
		EstimateCounts: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < t.N; i++ {
		for j := 0; j < 10000; j++ {
			data := make([]byte, 150)
			metadata := make([]byte, 50)
			rand.Read(data)
			rand.Read(metadata)

			var events []store.EventData
			events = append(events, store.EventData{
				Type:     "AccountOpened",
				Data:     data,
				Metadata: metadata,
			})

			_, err := eventstore.Add(uuid.New(), 0, events)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestAdd(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	eventstore, err := store.NewBadgerEventStore(store.BadgerStoreOptions{
		DB:             db,
		EstimateCounts: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	stream := uuid.New()

	data := make([]byte, 100)
	metadata := make([]byte, 50)
	rand.Read(data)
	rand.Read(metadata)

	events, err := eventstore.Add(stream, 0, []store.EventData{
		{
			Type:     "AccountOpened",
			Data:     data,
			Metadata: metadata,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.Equal(len(events), 1)
	assert.Equal(data, events[0].Data)
	assert.Equal(metadata, events[0].Metadata)
	assert.Equal(events[0].Stream, stream)
	assert.Equal(events[0].Type, "AccountOpened")
	assert.Equal(events[0].AddedAt.IsZero(), false)
	assert.NotEqual(events[0].ID, ulid.ULID{})
	assert.Equal(events[0].ID, events[0].CausationID)
	assert.Equal(events[0].ID, events[0].CorrelationID)
	assert.Equal(events[0].Version, uint32(0))
}

func TestGet(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	eventstore, err := store.NewBadgerEventStore(store.BadgerStoreOptions{
		DB:             db,
		EstimateCounts: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	stream := uuid.New()

	data := make([]byte, 100)
	metadata := make([]byte, 50)
	rand.Read(data)
	rand.Read(metadata)

	_, err = eventstore.Add(stream, 0, []store.EventData{
		{
			Type:     "AccountOpened",
			Data:     data,
			Metadata: metadata,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	events, err := eventstore.Get(stream, 0, 10)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.Equal(len(events), 1)
	assert.Equal(data, events[0].Data)
	assert.Equal(metadata, events[0].Metadata)
	assert.Equal(events[0].Stream, stream)
	assert.Equal(events[0].Type, "AccountOpened")
	assert.Equal(events[0].AddedAt.IsZero(), false)
	assert.NotEqual(events[0].ID, ulid.ULID{})
	assert.Equal(events[0].ID, events[0].CausationID)
	assert.Equal(events[0].ID, events[0].CorrelationID)
	assert.Equal(events[0].Version, uint32(0))
}

func TestGetAll(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	eventstore, err := store.NewBadgerEventStore(store.BadgerStoreOptions{
		DB:             db,
		EstimateCounts: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	stream := uuid.New()

	data := make([]byte, 100)
	metadata := make([]byte, 50)
	rand.Read(data)
	rand.Read(metadata)

	_, err = eventstore.Add(stream, 0, []store.EventData{
		{
			Type:     "AccountOpened",
			Data:     data,
			Metadata: metadata,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	events, err := eventstore.GetAll(ulid.ULID{}, 10)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.Equal(len(events), 1)
	assert.Equal(data, events[0].Data)
	assert.Equal(metadata, events[0].Metadata)
	assert.Equal(events[0].Stream, stream)
	assert.Equal(events[0].Type, "AccountOpened")
	assert.Equal(events[0].AddedAt.IsZero(), false)
	assert.NotEqual(events[0].ID, ulid.ULID{})
	assert.Equal(events[0].ID, events[0].CausationID)
	assert.Equal(events[0].ID, events[0].CorrelationID)
	assert.Equal(events[0].Version, uint32(0))
}
