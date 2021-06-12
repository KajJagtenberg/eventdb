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

	assert := assert.New(t)
	assert.Equal(len(res.Events), 1)
}
