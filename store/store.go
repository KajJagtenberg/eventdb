package store

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"math/rand"

	"github.com/oklog/ulid"
	"go.etcd.io/bbolt"
)

var (
	ErrConcurrentStreamModification = errors.New("Concurrent stream modification")
	ErrVersionNegative              = errors.New("Version cannot be negative")
	ErrEmptyEventData               = errors.New("List of event data is empty")
	ErrInvalidEventFormat           = errors.New("Invalid event format")
)

type EventStore struct {
	db      *bbolt.DB
	entropy io.Reader
}

func (store *EventStore) GetStreams(skip int, limit int) ([]Stream, error) {
	if skip < 0 {
		skip = 0
	}

	if limit < 0 {
		limit = 0
	}

	var streams []Stream

	if err := store.db.View(func(t *bbolt.Tx) error {
		cur := t.Bucket([]byte("streams")).Cursor()

		for k, v := cur.First(); k != nil && (limit == 0 || len(streams) < limit); k, v = cur.Next() {
			var stream Stream

			if err := json.Unmarshal(v, &stream); err != nil {
				return err
			}

			streams = append(streams, stream)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return streams, nil
}

func (store *EventStore) GetTotalStreams() (int, error) {
	var total int

	err := store.db.View(func(t *bbolt.Tx) error {
		stats := t.Bucket([]byte("stats"))

		v := stats.Get([]byte("total_streams"))

		if v != nil {
			total = int(binary.LittleEndian.Uint32(v))
		}
		return nil
	})

	return total, err
}

func (store *EventStore) GetTotalEvents() (int, error) {
	var total int

	err := store.db.View(func(t *bbolt.Tx) error {
		stats := t.Bucket([]byte("stats"))

		v := stats.Get([]byte("total_events"))

		if v != nil {
			total = int(binary.LittleEndian.Uint32(v))
		}

		return nil
	})

	return total, err
}

func NewEventStore(db *bbolt.DB) (*EventStore, error) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(int64(ulid.Now()))), 0)

	if err := db.Update(func(txn *bbolt.Tx) error {
		if _, err := txn.CreateBucketIfNotExists([]byte("streams")); err != nil {
			return err
		}

		if _, err := txn.CreateBucketIfNotExists([]byte("events")); err != nil {
			return err
		}

		if _, err := txn.CreateBucketIfNotExists([]byte("stats")); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &EventStore{db, entropy}, nil
}
