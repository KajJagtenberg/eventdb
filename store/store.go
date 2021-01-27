package store

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"github.com/vmihailenco/msgpack/v5"
	"go.etcd.io/bbolt"
)

var (
	entropy = ulid.Monotonic(rand.New((rand.NewSource((int64(ulid.Now()))))), 0)
)

var (
	ErrConcurrentStreamModifcation = errors.New("Concurrent stream modification")
)

type Store struct {
	db *bbolt.DB
}

func (s *Store) AppendToStream(streamId uuid.UUID, version int, events []AppendEvent) error {
	return s.db.Batch(func(txn *bbolt.Tx) error {
		streamsBucket := txn.Bucket([]byte("streams"))
		eventsBucket := txn.Bucket([]byte("events"))

		streamBucket := streamsBucket.Bucket(streamId[:])

		if streamBucket == nil {
			var err error

			streamBucket, err = streamsBucket.CreateBucket(streamId[:])

			if err != nil {
				return err
			}

			streamsBucket.NextSequence()
		}

		currentVersion := int(streamBucket.Sequence())

		if currentVersion != version {
			return ErrConcurrentStreamModifcation
		}

		for i, event := range events {
			id, err := ulid.New(ulid.Now(), entropy)
			if err != nil {
				return err
			}

			data, err := json.Marshal(event.Data)
			if err != nil {
				return err
			}

			serialized, err := msgpack.Marshal(Event{
				ID:            id,
				Stream:        streamId,
				Version:       version + i,
				Type:          event.Type,
				Data:          data,
				Metadata:      event.Metadata, // TODO: Perhaps turn this into json.RawMessage like data as well
				CausationID:   event.CausationID,
				CorrelationID: event.CorrelationID,
				Timestamp:     time.Now(),
			})
			if err != nil {
				return err
			}

			if err := eventsBucket.Put(id[:], serialized); err != nil {
				return err
			}

			if err := streamBucket.Put(itob(streamBucket.Sequence()), id[:]); err != nil {
				return err
			}

			streamBucket.NextSequence()
			eventsBucket.NextSequence()
		}

		return nil
	})
}

func (s *Store) LoadFromStream(streamId uuid.UUID, version int, limit int) ([]Event, int, error) {
	result := []Event{}
	total := 0

	err := s.db.View(func(txn *bbolt.Tx) error {
		streamsBucket := txn.Bucket([]byte("streams"))
		eventsBucket := txn.Bucket([]byte("events"))

		streamBucket := streamsBucket.Bucket(streamId[:])

		if streamBucket == nil {
			return nil
		}

		total = int(streamBucket.Sequence())

		cur := streamBucket.Cursor()

		for k, v := cur.First(); k != nil; k, v = cur.Next() {
			if version > 0 {
				version--
				continue
			}

			serialized := eventsBucket.Get(v)

			var event Event

			if err := msgpack.Unmarshal(serialized, &event); err != nil {
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

func (s *Store) Subscribe(offset ulid.ULID, limit int) ([]Event, error) {
	if limit == 0 {
		limit = 100
	}

	result := []Event{}

	err := s.db.View(func(txn *bbolt.Tx) error {
		cur := txn.Bucket([]byte("events")).Cursor()

		for k, v := cur.Seek(offset[:]); k != nil; k, v = cur.Next() {
			serialized := v

			var event Event

			if err := msgpack.Unmarshal(serialized, &event); err != nil {
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

func (s *Store) GetEventByID(id ulid.ULID) (Event, error) {
	var result Event

	err := s.db.View(func(txn *bbolt.Tx) error {
		bucket := txn.Bucket([]byte("events"))

		value := bucket.Get(id[:])

		if err := msgpack.Unmarshal(value, &result); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *Store) GetStreams(offset int, limit int) ([]uuid.UUID, int, error) {
	if limit == 0 {
		limit = 10
	}

	streams := []uuid.UUID{}
	total := 0

	err := s.db.View(func(txn *bbolt.Tx) error {
		bucket := txn.Bucket([]byte("streams"))

		total = int(bucket.Sequence())

		cur := bucket.Cursor()

		for k, _ := cur.First(); k != nil && len(streams) < limit; k, _ = cur.Next() {
			if offset > 0 {
				offset--
				continue
			}

			stream, err := uuid.FromBytes(k)
			if err != nil {
				return err
			}
			streams = append(streams, stream)
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return streams, total, nil
}

func (s *Store) GetEventCount() (int, error) {
	total := 0

	err := s.db.View(func(txn *bbolt.Tx) error {
		bucket := txn.Bucket([]byte("events"))

		total = int(bucket.Sequence())

		return nil
	})

	if err != nil {
		return 0, nil
	}

	return total, nil
}

func (s *Store) GetStreamCount() (int, error) {
	total := 0

	err := s.db.View(func(txn *bbolt.Tx) error {
		bucket := txn.Bucket([]byte("streams"))

		total = int(bucket.Sequence())

		return nil
	})

	if err != nil {
		return 0, nil
	}

	return total, nil
}

func (s *Store) GetDBSize() int64 {
	var size int64

	s.db.View(func(txn *bbolt.Tx) error {
		size = txn.Size()

		return nil
	})

	return size
}

func (s *Store) Backup(dst io.Writer) error {
	return s.db.View(func(txn *bbolt.Tx) error {
		_, err := txn.WriteTo(dst)

		return err
	})
}

func NewStore(db *bbolt.DB) (*Store, error) {
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

	return &Store{db}, nil
}

func itob(v uint64) []byte {
	r := make([]byte, 8)
	binary.BigEndian.PutUint64(r, v)
	return r
}
