package store_test

import (
	"eventflowdb/store"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"github.com/stretchr/testify/assert"
)

func TestSerializeRecordedEvent(t *testing.T) {
	event := store.RecordedEvent{
		ID:      ulid.MustNew(ulid.Now(), rand.New(rand.NewSource(int64(ulid.Now())))),
		Stream:  uuid.New(),
		Version: 10,
		Type:    "ProductAdded",
		Data:    []byte(`{"name":"Appeltaart","price":400}`),
		AddedAt: time.Now(),
	}

	serialized, err := event.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(serialized)
	t.Log(len(serialized))
	t.Log(len(serialized) * 8)
}

func TestDeserializeRecordedEvent(t *testing.T) {
	entropy := rand.New(rand.NewSource(int64(ulid.Now())))

	assert := assert.New(t)
	event := store.RecordedEvent{
		ID:            ulid.MustNew(ulid.Now(), entropy),
		Stream:        uuid.New(),
		Version:       10,
		Type:          "ProductAdded",
		Data:          []byte(`{"name":"Appeltaart","price":400}`),
		AddedAt:       time.Now(),
		CausationID:   ulid.MustNew(ulid.Now(), entropy),
		CorrelationID: ulid.MustNew(ulid.Now(), entropy),
	}

	serialized, err := event.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	deserialized := store.RecordedEvent{}
	if err := deserialized.Deserialize(serialized); err != nil {
		t.Fatal(err)
	}

	assert.Equal(event.ID, deserialized.ID)
	assert.Equal(event.Stream, deserialized.Stream)
	assert.Equal(event.Version, deserialized.Version)
	assert.Equal(event.AddedAt.Equal(deserialized.AddedAt), true)
	assert.Equal(event.Type, deserialized.Type)
	assert.Equal(event.Data, deserialized.Data)

	t.Log(deserialized)
}
