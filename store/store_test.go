package store

import (
	"encoding/binary"
	"math/rand"
	"os"
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/bbolt"
)

func TempStore() (EventStore, error) {
	path := "tmp/badger/*"
	os.RemoveAll(path)

	db, err := badger.Open(badger.DefaultOptions(path).WithLogger(nil))
	if err != nil {
		return nil, err
	}

	return NewBadgerEventStore(BadgerStoreOptions{
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

func TestEventCount(t *testing.T) {
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

	req := &api.EventCountRequest{}

	res, err := store.EventCount(req)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.Equal(int64(1), res.Count)
}

func TestStreamCount(t *testing.T) {
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

	req := &api.StreamCountRequest{}

	res, err := store.StreamCount(req)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.Equal(int64(1), res.Count)
}

func BenchmarkAdd(t *testing.B) {
	store, err := TempStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	for i := 0; i < t.N; i++ {
		var data []*api.AddRequest_EventData

		for j := 0; j < 1; j++ {
			data = append(data, &api.AddRequest_EventData{
				Type:     "TestEvent",
				Data:     []byte("data"),
				Metadata: []byte("metadata"),
			})
		}

		req := &api.AddRequest{
			Stream:  uuid.New().String(),
			Version: 0,
			Events:  data,
		}

		if _, err := store.Add(req); err != nil {
			t.Fatal(err)
		}
	}
}

func BenchmarkTxn(b *testing.B) {
	db, err := bbolt.Open("tmp/store.db", 0666, bbolt.DefaultOptions)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < b.N; i++ {
		txn, err := db.Begin(true)
		if err != nil {
			b.Fatal(err)
		}
		defer txn.Rollback()

		bucket, err := txn.CreateBucketIfNotExists([]byte("bucket"))
		if err != nil {
			b.Fatal(err)
		}

		key := make([]byte, 16)
		value := make([]byte, 150)
		rand.Read(key)
		rand.Read(value)

		if err := bucket.Put(key, value); err != nil {
			b.Fatal(err)
		}

		if err := txn.Commit(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBadgerWrites(t *testing.B) {
	db, err := badger.Open(badger.DefaultOptions("tmp/badger").WithLogger(nil))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < t.N; i++ {
		txn := db.NewTransaction(true)
		defer txn.Discard()

		key := make([]byte, 4)
		value := make([]byte, 150)

		rand.Read(value)

		binary.BigEndian.PutUint32(key, uint32(i))

		if err := txn.Set(key, value); err != nil {
			t.Fatal(err)
		}

		if err := txn.Commit(); err != nil {
			t.Fatal(err)
		}
	}
}
