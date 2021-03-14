package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"time"

	"github.com/google/uuid"
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

func (store *EventStore) AppendToStream(stream uuid.UUID, version int, events []EventData) ([]RecordedEvent, error) {
	var records []RecordedEvent

	err := store.db.Update(func(t *bbolt.Tx) error {
		streamsBucket := t.Bucket([]byte("streams"))
		eventsBucket := t.Bucket([]byte("events"))

		var persistedStream Stream

		v := streamsBucket.Get(stream[:])

		if v != nil {
			if err := json.Unmarshal(v, &persistedStream); err != nil {
				return err
			}
		}

		if version != persistedStream.Size() {
			return ErrConcurrentStreamModification
		}

		for i, event := range events {
			id, err := ulid.New(ulid.Now(), store.entropy)
			if err != nil {
				return err
			}

			record := RecordedEvent{
				ID:       id,
				Version:  version + i,
				Stream:   stream,
				Type:     event.Type,
				Data:     event.Data,
				Metadata: event.Metadata,
				AddedAt:  time.Now(),
			}

			records = append(records, record)

			persistedStream.Events = append(persistedStream.Events, record.ID)

			v, err := json.Marshal(record)
			if err != nil {
				return err
			}

			if err := eventsBucket.Put(record.ID[:], v); err != nil {
				return err
			}
		}

		v, err := json.Marshal(persistedStream)
		if err != nil {
			return err
		}

		if err := streamsBucket.Put(stream[:], v); err != nil {
			return err
		}

		return nil
	})

	return records, err
}

func (store *EventStore) LoadFromAll(offset ulid.ULID, limit int) ([]RecordedEvent, error) {
	if limit < 0 {
		limit = 0
	}

	var records []RecordedEvent

	err := store.db.View(func(t *bbolt.Tx) error {
		cur := t.Bucket([]byte("events")).Cursor()

		k, _ := cur.Seek(offset[:])

		if bytes.Compare(k, offset[:]) == 0 {
			cur.Next()
		}

		for k, v := cur.Seek(offset[:]); k != nil && (limit == 0 || len(records) < limit); k, v = cur.Next() {
			if bytes.Compare(k, offset[:]) == 0 {
				continue
			}

			var record RecordedEvent

			if err := json.Unmarshal(v, &record); err != nil {
				return err
			}

			records = append(records, record)
		}

		return nil
	})

	return records, err
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
		cur := t.Bucket([]byte("streams")).Cursor()

		for k, _ := cur.First(); k != nil; k, _ = cur.Next() {
			total++
		}

		return nil
	})

	return total, err
}

func (store *EventStore) GetTotalEvents() (int, error) {
	var total int

	err := store.db.View(func(t *bbolt.Tx) error {
		cur := t.Bucket([]byte("events")).Cursor()

		for k, _ := cur.First(); k != nil; k, _ = cur.Next() {
			total++
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

		return nil
	}); err != nil {
		return nil, err
	}

	return &EventStore{db, entropy}, nil
}
