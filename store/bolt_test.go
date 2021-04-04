package store

import (
	"io/ioutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/bbolt"
)

func OpenStore() (*BoltStore, error) {
	f, err := ioutil.TempFile("/tmp", "*")
	if err != nil {
		return nil, err
	}

	db, err := bbolt.Open(f.Name(), 0666, bbolt.DefaultOptions)
	if err != nil {
		return nil, err
	}

	store, err := NewBoltStore(db)
	if err != nil {
		return nil, err
	}
	return store, nil
}

func TestBoltAdd(t *testing.T) {
	assert := assert.New(t)

	store, err := OpenStore()
	if err != nil {
		t.Fatal(err)
	}

	stream := uuid.New()
	events := []EventData{
		{
			Type:     "TestEvent",
			Data:     []byte("data"),
			Metadata: []byte("metadata"),
		},
	}

	records, err := store.Add(stream, 0, events)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(len(records), len(events))

	for i := 0; i < len(events); i++ {
		assert.Equal(events[i].Type, records[i].Type)
		assert.Equal(events[i].Data, records[i].Data)
		assert.Equal(events[i].Metadata, records[i].Metadata)
		assert.NotEqual(events[i].CausationID, records[i].CausationID)
		assert.NotEqual(events[i].CorrelationID, records[i].CorrelationID)
	}
}

func TestBoltAddConcurrency(t *testing.T) {
	assert := assert.New(t)

	store, err := OpenStore()
	if err != nil {
		t.Fatal(err)
	}

	stream := uuid.New()
	events := []EventData{
		{
			Type:     "TestEvent",
			Data:     []byte("data"),
			Metadata: []byte("metadata"),
		},
	}

	_, err = store.Add(stream, 0, events)
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.Add(stream, 0, events)
	assert.NotEqual(err, nil)
}

func TestBoltGet(t *testing.T) {
	assert := assert.New(t)

	store, err := OpenStore()
	if err != nil {
		t.Fatal(err)
	}

	stream := uuid.New()
	events := []EventData{
		{
			Type:     "TestEvent",
			Data:     []byte("data"),
			Metadata: []byte("metadata"),
		},
	}

	_, err = store.Add(stream, 0, events)
	if err != nil {
		t.Fatal(err)
	}

	records, err := store.Get(stream, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(len(records), len(events))

	for i := 0; i < len(events); i++ {
		assert.Equal(events[i].Type, records[i].Type)
		assert.Equal(events[i].Data, records[i].Data)
		assert.Equal(events[i].Metadata, records[i].Metadata)
		assert.NotEqual(events[i].CausationID, records[i].CausationID)
		assert.NotEqual(events[i].CorrelationID, records[i].CorrelationID)
	}
}
