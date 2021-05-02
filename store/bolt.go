package store

import (
	"bytes"
	"errors"
	"hash/crc32"
	"io"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
)

var (
	buckets = []string{"streams", "events"}
	entropy = ulid.Monotonic(rand.New(rand.NewSource(int64(ulid.Now()))), 0)
)

const (
	ESTIMATE_SLEEP_TIME = time.Second // TODO: Maybe make this configurable?
)

type BoltStore struct {
	db                  *bbolt.DB
	estimateStreamCount int64
	estimateEventCount  int64
}

func (s *BoltStore) Size() (int64, error) {
	var size int64 = 0

	err := s.db.View(func(t *bbolt.Tx) error {
		size = t.Size()

		return nil
	})

	return size, err
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
	if bytes.Equal(stream[:], make([]byte, 16)) {
		return nil, errors.New("stream cannot be all zeroes")
	}

	if len(events) == 0 {
		return nil, errors.New("list of events is empty")
	}

	result := make([]Event, 0)

	if err := s.db.Batch(func(t *bbolt.Tx) error {
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

		if int(version) < len(s.Events) {
			return errors.New("concurrent stream modification")
		}

		if int(version) > len(s.Events) {
			return errors.New("given version leaves gap in stream")
		}

		now := time.Now()

		for i, event := range events {
			if event.Type == "" {
				return errors.New("event type cannot be empty")
			}

			if len(event.Data) == 0 {
				return errors.New("event data cannot be empty")
			}

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
				AddedAt:       now,
			}

			if bytes.Equal(record.CausationID[:], make([]byte, 16)) {
				record.CausationID = record.ID
			}

			if bytes.Equal(record.CorrelationID[:], make([]byte, 16)) {
				record.CorrelationID = record.CausationID
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
	result := make([]Event, 0)

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
				return errors.New("event cannot be found. this should never happen")
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *BoltStore) GetAll(offset ulid.ULID, limit uint32) ([]Event, error) {
	if limit == 0 {
		limit = 100
	}

	result := make([]Event, 0)

	if err := s.db.View(func(t *bbolt.Tx) error {
		cursor := t.Bucket([]byte("events")).Cursor()

		for k, v := cursor.Seek(offset[:]); k != nil; k, v = cursor.Next() {
			if len(result) >= int(limit) {
				break
			}

			var event Event
			if err := event.Unmarshal(v); err != nil {
				return err
			}

			result = append(result, event)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *BoltStore) EventCount() (int64, error) {
	var total int64

	if err := s.db.View(func(t *bbolt.Tx) error {
		cursor := t.Bucket([]byte("events")).Cursor()

		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			total++
		}

		return nil
	}); err != nil {
		return 0, err
	}

	return total, nil
}

func (s *BoltStore) StreamCount() (int64, error) {
	var total int64

	if err := s.db.View(func(t *bbolt.Tx) error {
		cursor := t.Bucket([]byte("streams")).Cursor()

		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			total++
		}

		return nil
	}); err != nil {
		return 0, err
	}

	return total, nil
}
func (s *BoltStore) StreamCountEstimate() (int64, error) {
	return s.estimateStreamCount, nil
}

func (s *BoltStore) EventCountEstimate() (int64, error) {
	return s.estimateEventCount, nil
}

// TODO: Store the checksum and ID at intervals to prevent recalculation since the beginning
func (s *BoltStore) Checksum() (id ulid.ULID, sum []byte, err error) {
	crc := crc32.NewIEEE()

	err = s.db.View(func(t *bbolt.Tx) error {
		cursor := t.Bucket([]byte("events")).Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if err := id.UnmarshalBinary(k); err != nil {
				return err
			}

			if _, err := crc.Write(v); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return id, sum, err
	}

	sum = crc.Sum(nil)

	return id, sum, nil
}

func (s *BoltStore) Close() error {
	return s.db.Close()
}

func NewBoltStore(db *bbolt.DB, log logrus.StdLogger) (*BoltStore, error) {
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

	store := &BoltStore{db, 0, 0}

	go func() {
		for {
			streamCount, err := store.StreamCount()
			if err != nil {
				log.Fatalf("Failed to get stream count: %v", err)
			}

			eventCount, err := store.EventCount()
			if err != nil {
				log.Fatalf("Failed to get stream count: %v", err)
			}

			store.estimateStreamCount = streamCount
			store.estimateEventCount = eventCount

			time.Sleep(ESTIMATE_SLEEP_TIME)
		}
	}()

	return store, nil
}
