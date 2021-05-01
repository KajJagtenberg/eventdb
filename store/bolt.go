package store

import (
	"bytes"
	"errors"
	"hash/crc32"
	"io"
	"log"
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

const (
	ESTIMATE_SLEEP_TIME = time.Second // Maybe make this configurable?
)

// var (
// 	addCounter = promauto.NewCounter(prometheus.CounterOpts{
// 		Name: "store_add_total",
// 		Help: "The amount of events that have been added to the store",
// 	})
// 	getCounter = promauto.NewCounter(prometheus.CounterOpts{
// 		Name: "store_get_total",
// 		Help: "The amount of get requests performed",
// 	})
// 	logCounter = promauto.NewCounter(prometheus.CounterOpts{
// 		Name: "store_log_total",
// 		Help: "The amount of log requests performed",
// 	})
// 	concurrencyCounter = promauto.NewCounter(prometheus.CounterOpts{
// 		Name: "store_concurrent_modification_total",
// 		Help: "The total amount of concurrent stream modification errors",
// 	})
// )

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
	if bytes.Compare(stream[:], make([]byte, 16)) == 0 {
		return nil, errors.New("Stream cannot be all zeroes")
	}

	if version < 0 {
		log.Printf("Version is negative")
		return nil, errors.New("Version cannot be negative")
	}

	if len(events) == 0 {
		return nil, errors.New("List of events is empty")
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
			return errors.New("Concurrent stream modification")
		}

		if int(version) > len(s.Events) {
			return errors.New("Given version leaves gap in stream")
		}

		now := time.Now()

		for i, event := range events {
			if event.Type == "" {
				return errors.New("Event type cannot be empty")
			}

			if len(event.Data) == 0 {
				return errors.New("Event data cannot be empty")
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

	// addCounter.Add(float64(len(events)))

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
				return errors.New("Event cannot be found. This should never happen.")
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// getCounter.Add(1)

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

	// logCounter.Add(1)

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

func (s *BoltStore) Checksum() (uint32, error) {
	crc := crc32.NewIEEE()

	err := s.db.View(func(t *bbolt.Tx) error {
		cursor := t.Bucket([]byte("events")).Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if _, err := crc.Write(v); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return 0, nil
	}

	return crc.Sum32(), nil
}

func (s *BoltStore) Close() error {
	return s.db.Close()
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
