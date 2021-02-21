package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"time"

	"github.com/go-playground/validator/v10"
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

func (store *EventStore) AppendToStream(name uuid.UUID, version int, events []EventData) ([]RecordedEvent, error) {
	if version < 0 {
		return nil, ErrVersionNegative
	}

	if len(events) == 0 {
		return nil, ErrEmptyEventData
	}

	var result []RecordedEvent

	err := store.db.Update(func(txn *bbolt.Tx) error {
		streamBucket := txn.Bucket([]byte("streams"))
		eventBucket := txn.Bucket([]byte("events"))

		var stream Stream
		stream.Unmarshal(streamBucket.Get(name[:]))

		if len(stream.Events) != version {
			return ErrConcurrentStreamModification
		}

		validator := validator.New()

		for i, event := range events {
			if err := validator.Struct(event); err != nil {
				return ErrInvalidEventFormat
			}

			id := ulid.MustNew(ulid.Now(), store.entropy)

			recorded := RecordedEvent{
				ID:       id,
				Stream:   name,
				Version:  version + i,
				Type:     event.Type,
				Data:     event.Data,
				Metadata: event.Metadata,
				AddedAt:  time.Now(),
			}

			result = append(result, recorded)

			serialized, err := json.Marshal(recorded)
			if err != nil {
				return err
			}

			if err := eventBucket.Put(id[:], serialized); err != nil {
				return err
			}

			stream.Events = append(stream.Events, id)
		}

		if err := streamBucket.Put(name[:], stream.Marshal()); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (store *EventStore) LoadFromStream(name uuid.UUID, version int, limit int) ([]RecordedEvent, int, error) {
	var result []RecordedEvent
	var total int

	err := store.db.View(func(txn *bbolt.Tx) error {
		streamBucket := txn.Bucket([]byte("streams"))
		eventBucket := txn.Bucket([]byte("events"))

		var stream Stream
		stream.Unmarshal(streamBucket.Get(name[:]))

		total = len(stream.Events)

		for _, id := range stream.Events {
			if version > 0 {
				continue
			}

			if len(result) == limit && limit != 0 {
				break
			}

			serialized := eventBucket.Get(id[:])

			var event RecordedEvent

			if err := json.Unmarshal(serialized, &event); err != nil {
				return err
			}

			result = append(result, event)
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (store *EventStore) LoadFromAll(offset ulid.ULID, limit int) ([]RecordedEvent, error) {
	var result []RecordedEvent

	err := store.db.View(func(txn *bbolt.Tx) error {
		cursor := txn.Bucket([]byte("events")).Cursor()

		for k, v := cursor.Seek(offset[:]); k != nil && (len(result) < limit || limit == 0); k, v = cursor.Next() {
			if bytes.Compare(offset[:], k) == 0 {
				continue
			}

			var event RecordedEvent

			if err := json.Unmarshal(v, &event); err != nil {
				return err
			}

			result = append(result, event)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (store *EventStore) GetStreams(offset int, limit int) ([]uuid.UUID, int, error) {
	var result []uuid.UUID
	var total int

	if offset < 0 {
		offset = 0
	}

	if limit < 0 {
		limit = 0
	}

	err := store.db.View(func(txn *bbolt.Tx) error {
		cursor := txn.Bucket([]byte("streams")).Cursor()

		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			total++

			if offset > 0 {
				continue
			}

			if len(result) == limit && limit != 0 {
				continue
			}

			id, err := uuid.FromBytes(k)
			if err != nil {
				return err
			}

			result = append(result, id)
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (store *EventStore) GetEventByID(id ulid.ULID) (RecordedEvent, error) {
	var event RecordedEvent

	err := store.db.View(func(txn *bbolt.Tx) error {
		if value := txn.Bucket([]byte("events")).Get(id[:]); value != nil {
			return json.Unmarshal(value, &event)
		}

		return nil
	})

	return event, err
}

func (store *EventStore) GetDBSize() int64 {
	var size int64

	store.db.View(func(txn *bbolt.Tx) error {
		size = txn.Size()

		return nil
	})

	return size
}

func (store *EventStore) GetEventCount() (int, error) {
	var count int

	err := store.db.View(func(txn *bbolt.Tx) error {
		cursor := txn.Bucket([]byte("events")).Cursor()

		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			count++
		}

		return nil
	})

	return count, err
}

func (store *EventStore) GetStreamCount() (int, error) {
	var count int

	err := store.db.View(func(txn *bbolt.Tx) error {
		cursor := txn.Bucket([]byte("streams")).Cursor()

		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			count++
		}

		return nil
	})

	return count, err
}

func (store *EventStore) Backup(dst io.Writer) error {
	return store.db.View(func(txn *bbolt.Tx) error {
		_, err := txn.WriteTo(dst)

		return err
	})
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
