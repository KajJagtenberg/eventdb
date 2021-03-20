package store_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/bbolt"
)

func TestAdd(t *testing.T) {
	assert := assert.New(t)

	db, err := bbolt.Open("/tmp/events.db", 0666, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	storage, err := store.NewStorage(db)
	if err != nil {
		t.Fatal(err)
	}

	stream, err := uuid.New().MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	req := &store.AddRequest{
		Stream:  stream,
		Version: 0,
		Events: []*store.AddRequest_Event{
			{
				Type: "ProductAdded",
				Data: []byte(`{"name":"Samsung Galaxy S8","version":80000}`),
			},
		},
	}

	records, err := storage.Add(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(len(records), 1)
	assert.Equal(records[0].Stream, req.Stream)
	assert.Equal(records[0].Type, req.Events[0].Type)
	assert.Equal(records[0].Data, req.Events[0].Data)
}

func TestGet(t *testing.T) {
	assert := assert.New(t)

	db, err := bbolt.Open("/tmp/events.db", 0666, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	storage, err := store.NewStorage(db)
	if err != nil {
		t.Fatal(err)
	}

	stream, err := uuid.New().MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	req := &store.AddRequest{
		Stream:  stream,
		Version: 0,
		Events: []*store.AddRequest_Event{
			{
				Type: "ProductAdded",
				Data: []byte(`{"name":"Samsung Galaxy S8","version":80000}`),
			},
		},
	}

	records, err := storage.Add(req)
	if err != nil {
		t.Fatal(err)
	}

	records, err = storage.Get(&store.GetRequest{
		Stream:  req.Stream,
		Version: 0,
		Limit:   0,
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(len(records), 1)
	assert.Equal(records[0].Stream, req.Stream)
	assert.Equal(records[0].Type, req.Events[0].Type)
	assert.Equal(records[0].Data, req.Events[0].Data)
}

func TestLog(t *testing.T) {
	assert := assert.New(t)

	db, err := bbolt.Open("/tmp/events.db", 0666, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	storage, err := store.NewStorage(db)
	if err != nil {
		t.Fatal(err)
	}

	stream, err := uuid.New().MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	req := &store.AddRequest{
		Stream:  stream,
		Version: 0,
		Events: []*store.AddRequest_Event{
			{
				Type: "ProductAdded",
				Data: []byte(`{"name":"Samsung Galaxy S8","version":80000}`),
			},
		},
	}

	records, err := storage.Add(req)
	if err != nil {
		t.Fatal(err)
	}

	records, err = storage.Log(&store.LogRequest{
		Offset: make([]byte, 16),
		Limit:  0,
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(len(records) > 0, true)
}
