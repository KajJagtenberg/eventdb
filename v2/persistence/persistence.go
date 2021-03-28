package persistence

import (
	"bytes"
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"go.etcd.io/bbolt"
)

const (
	BUCKET_STREAMS = "streams"
	BUCKET_EVENTS  = "events"
)

var (
	entropy = ulid.Monotonic(rand.New(rand.NewSource(int64(ulid.Now()))), 0)
)

type Persistence struct {
	db *bbolt.DB
}

func (p *Persistence) Add(streamID uuid.UUID, version uint32, events []EventData) ([]Event, error) {
	var result []Event

	err := p.db.Update(func(t *bbolt.Tx) error {
		streamsBucket := t.Bucket([]byte(BUCKET_STREAMS))
		eventsBucket := t.Bucket([]byte(BUCKET_EVENTS))

		var stream Stream

		if packed := streamsBucket.Get(streamID[:]); packed != nil {
			if err := stream.Unmarshal(packed); err != nil {
				return err
			}
		} else {
			stream.AddedAt = time.Now()
			stream.ID = streamID
		}

		if len(stream.Events) != int(version) {
			return errors.New("Concurrent stream modification")
		}

		for i, event := range events {
			id, err := ulid.New(ulid.Now(), entropy)
			if err != nil {
				return err
			}

			record := Event{}
			record.ID = id
			record.Stream = streamID
			record.Version = version + uint32(i)
			record.Data = event.Data
			record.Metadata = event.Metadata
			record.CausationID = event.CausationID
			record.CorrelationID = event.CorrelationID

			if bytes.Compare(record.CausationID[:], make([]byte, 16)) == 0 {
				record.CausationID = record.ID
			}

			if bytes.Compare(record.CorrelationID[:], make([]byte, 16)) == 0 {
				record.CorrelationID = record.ID
			}

			record.AddedAt = time.Now()

			stream.Events = append(stream.Events, record.ID)

			packed, err := record.Marshal()
			if err != nil {
				return err
			}

			if err := eventsBucket.Put(record.ID[:], packed); err != nil {
				return err
			}
		}

		packed, err := stream.Marshal()
		if err != nil {
			return err
		}

		return streamsBucket.Put(stream.ID[:], packed)
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewPersistence(db *bbolt.DB) (*Persistence, error) {
	err := db.Update(func(t *bbolt.Tx) error {
		buckets := []string{BUCKET_STREAMS, BUCKET_EVENTS}

		for _, bucket := range buckets {
			if _, err := t.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Persistence{db}, nil
}

type EventData struct {
	Type          string
	Data          []byte
	Metadata      []byte
	CausationID   ulid.ULID
	CorrelationID ulid.ULID
	AddedAt       time.Time
}

type Event struct {
	ID            ulid.ULID
	Stream        uuid.UUID
	Version       uint32
	Type          string
	Data          []byte
	Metadata      []byte
	CausationID   ulid.ULID
	CorrelationID ulid.ULID
	AddedAt       time.Time
}

func (e *Event) Marshal() ([]byte, error) {
	return nil, errors.New("Event Marshal not implemented")
}

func (e *Event) Unmarshal(data []byte) error {
	return errors.New("Event Unmarshal not implemented")
}

type Stream struct {
	ID      uuid.UUID
	Events  []ulid.ULID
	AddedAt time.Time
}

func (s *Stream) Marshal() ([]byte, error) {
	return nil, errors.New("Stream Marshal not implemented")
}

func (s *Stream) Unmarshal(data []byte) error {
	return errors.New("Stream Unmarshal not implemented")
}
