package store

import (
	"io/ioutil"
	"testing"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/bbolt"
)

func TempStore() (EventStore, error) {
	f, err := ioutil.TempFile("/tmp", "*")
	if err != nil {
		return nil, err
	}

	db, err := bbolt.Open(f.Name(), 0666, bbolt.DefaultOptions)
	if err != nil {
		return nil, err
	}

	return NewBoltEventStore(BoltStoreOptions{
		DB: db,
	})
}

func TestAdd(t *testing.T) {
	store, err := TempStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	req := &api.AddRequest{
		Stream:  uuid.New().String(),
		Version: 0,
		Events: []*api.AddRequest_EventData{
			{
				Type:     "TestEvent",
				Data:     []byte("data"),
				Metadata: []byte("metadata"),
			},
		},
	}

	res, err := store.Add(req)
	if err != nil {
		t.Fatal(err)
	}

	events := res.Events

	assert := assert.New(t)
	assert.Equal(1, len(events))
	assert.Equal(26, len(events[0].Id))
	assert.Equal(26, len(events[0].CausationId))
	assert.Equal(26, len(events[0].CorrelationId))
	assert.Equal([]byte("data"), events[0].Data)
	assert.Equal([]byte("metadata"), events[0].Metadata)
	assert.Equal(events[0].Id, events[0].CausationId)
	assert.Equal(events[0].Id, events[0].CorrelationId)
}

func TestGet(t *testing.T) {
	store, err := TempStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	stream := uuid.New().String()

	func() {
		req := &api.AddRequest{
			Stream:  stream,
			Version: 0,
			Events: []*api.AddRequest_EventData{
				{
					Type:     "TestEvent",
					Data:     []byte("data"),
					Metadata: []byte("metadata"),
				},
			},
		}

		_, err := store.Add(req)
		if err != nil {
			t.Fatal(err)
		}
	}()

	req := &api.GetRequest{
		Stream:  stream,
		Version: 0,
		Limit:   10,
	}

	res, err := store.Get(req)
	if err != nil {
		t.Fatal(err)
	}

	events := res.Events

	assert := assert.New(t)
	assert.Equal(1, len(events))
	assert.Equal(26, len(events[0].Id))
	assert.Equal(26, len(events[0].CausationId))
	assert.Equal(26, len(events[0].CorrelationId))
	assert.Equal([]byte("data"), events[0].Data)
	assert.Equal([]byte("metadata"), events[0].Metadata)
	assert.Equal(events[0].Id, events[0].CausationId)
	assert.Equal(events[0].Id, events[0].CorrelationId)
}

func TestGetAll(t *testing.T) {
	store, err := TempStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	stream := uuid.New().String()

	func() {
		req := &api.AddRequest{
			Stream:  stream,
			Version: 0,
			Events: []*api.AddRequest_EventData{
				{
					Type:     "TestEvent",
					Data:     []byte("data"),
					Metadata: []byte("metadata"),
				},
			},
		}

		_, err := store.Add(req)
		if err != nil {
			t.Fatal(err)
		}
	}()

	req := &api.GetAllRequest{}

	res, err := store.GetAll(req)
	if err != nil {
		t.Fatal(err)
	}

	events := res.Events

	assert := assert.New(t)
	assert.Equal(1, len(events))
	assert.Equal(26, len(events[0].Id))
	assert.Equal(26, len(events[0].CausationId))
	assert.Equal(26, len(events[0].CorrelationId))
	assert.Equal([]byte("data"), events[0].Data)
	assert.Equal([]byte("metadata"), events[0].Metadata)
	assert.Equal(events[0].Id, events[0].CausationId)
	assert.Equal(events[0].Id, events[0].CorrelationId)
}
