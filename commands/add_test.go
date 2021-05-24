package commands_test

import (
	"crypto/rand"
	"encoding/json"
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/commands"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/kajjagtenberg/go-commando"
	"github.com/oklog/ulid"
	"github.com/stretchr/testify/assert"
)

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

	handler := commands.AddHandler(eventstore)

	stream := uuid.New()

	data := make([]byte, 100)
	metadata := make([]byte, 50)
	rand.Read(data)
	rand.Read(metadata)

	args, err := json.Marshal(commands.AddRequest{
		Stream:  stream,
		Version: 0,
		Events: []store.EventData{
			{
				Type:     "AccountOpened",
				Data:     data,
				Metadata: metadata,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := handler(commando.Command{
		Name: "add",
		Args: args,
	})
	if err != nil {
		t.Fatal(err)
	}

	res, ok := result.(commands.AddResponse)
	if !ok {
		t.Fatal("wrong cast")
	}

	events := res.Events

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
