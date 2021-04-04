package store

import (
	"bytes"
	"errors"
	"io"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"go.etcd.io/bbolt"
)

var (
	buckets = []string{"streams", "events"}
	entropy = ulid.Monotonic(rand.New(rand.NewSource(int64(ulid.Now()))), 0)
)

type BoltStore struct {
	db *bbolt.DB
}

func (s *BoltStore) Size() int64 {
	var size int64 = 0

	s.db.View(func(t *bbolt.Tx) error {
		size = t.Size()

		return nil
	})

	return size
}

func (s *BoltStore) Backup(dst io.Writer) error {
	return s.db.View(func(t *bbolt.Tx) error {
		if _, err := t.WriteTo(dst); err != nil {
			return err
		}

		return nil
	})
}

func (s *BoltStore) Add(stream uuid.UUID, version uint32, events []EventData) ([]Event, error) {
	var result []Event

	if err := s.db.Update(func(t *bbolt.Tx) error {
		streamBucket := t.Bucket([]byte("streams"))
		eventsBucket := t.Bucket([]byte("events"))

		var s Stream

		if value := streamBucket.Get(stream[:]); value != nil {
			if err := s.Unmarshal(value); err != nil {
				return err
			}
		} else {
			s.ID = stream
			s.AddedAt = time.Now()
		}

		if len(s.Events) != int(version) {
			return errors.New("Concurrent stream modification")
		}

		now := time.Now()

		for i, event := range events {
			id, err := ulid.New(ulid.Now(), entropy)
			if err != nil {
				return err
			}

			record := Event{
				ID:            id,
				Stream:        stream,
				Version:       version + uint32(i),
				Type:          event.Type,
				Data:          event.Data,
				Metadata:      event.Metadata,
				CausationID:   event.CausationID,
				CorrelationID: event.CorrelationID,
				AddedAt:       event.AddedAt,
			}

			if bytes.Compare(record.CausationID[:], make([]byte, 16)) == 0 {
				record.CausationID = record.ID
			}

			if bytes.Compare(record.CorrelationID[:], make([]byte, 16)) == 0 {
				record.CorrelationID = record.CausationID
			}

			if record.AddedAt.IsZero() {
				record.AddedAt = now
			}

			result = append(result, record)

			s.Events = append(s.Events, record.ID)
		}

		for _, record := range result {
			if value, err := record.Marshal(); err != nil {
				return err
			} else {
				if err := eventsBucket.Put(record.ID[:], value); err != nil {
					return err
				}
			}
		}

		if value, err := s.Marshal(); err != nil {
			return err
		} else {
			if err := streamBucket.Put(s.ID[:], value); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *BoltStore) Get(stream uuid.UUID, version uint32, limit uint32) ([]Event, error) {
	var result []Event

	if err := s.db.View(func(t *bbolt.Tx) error {
		streamBucket := t.Bucket([]byte("streams"))
		eventsBucket := t.Bucket([]byte("events"))

		var s Stream

		if value := streamBucket.Get(stream[:]); value != nil {
			if err := s.Unmarshal(value); err != nil {
				return err
			}
		} else {
			return nil
		}

		for _, id := range s.Events {
			if version > 0 {
				version--
				continue
			}
			if len(result) >= int(limit) && limit != 0 {
				return nil
			}

			if value := eventsBucket.Get(id[:]); value != nil {
				var event Event
				if err := event.Unmarshal(value); err != nil {
					return err
				}

				result = append(result, event)
			} else {
				return errors.New("Event cannot be found. This should never happen.")
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func NewBoltStore(db *bbolt.DB) (*BoltStore, error) {
	if err := db.Update(func(t *bbolt.Tx) error {
		for _, bucket := range buckets {
			if _, err := t.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &BoltStore{db}, nil
}
