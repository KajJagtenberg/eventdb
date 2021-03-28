package persistence

import (
	"bytes"
	"errors"
	"math/rand"
	"time"

	proto "github.com/golang/protobuf/proto"
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

			result = append(result, record)
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

func (p *Persistence) Get(streamID uuid.UUID, version uint32, limit uint32) ([]Event, error) {
	var result []Event

	err := p.db.View(func(t *bbolt.Tx) error {
		streamsBucket := t.Bucket([]byte(BUCKET_STREAMS))
		eventsBucket := t.Bucket([]byte(BUCKET_EVENTS))

		var stream Stream

		if packed := streamsBucket.Get(streamID[:]); packed != nil {
			if err := stream.Unmarshal(packed); err != nil {
				return err
			}
		} else {
			return nil
		}

		for _, id := range stream.Events {
			if version > 0 {
				version--
				continue
			}
			if version != 0 && len(result) >= int(limit) {
				break
			}

			packed := eventsBucket.Get(id[:])

			record := Event{}
			if err := record.Unmarshal(packed); err != nil {
				return err
			}

			result = append(result, record)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Persistence) Log(offset ulid.ULID, limit uint32) ([]Event, error) {
	if limit == 0 {
		limit = 10
	}

	var result []Event

	err := p.db.View(func(t *bbolt.Tx) error {
		cursor := t.Bucket([]byte(BUCKET_EVENTS)).Cursor()

		for k, v := cursor.Seek(offset[:]); k != nil && len(result) < int(limit); k, v = cursor.Next() {
			if bytes.Compare(k, offset[:]) == 0 {
				continue
			}

			record := Event{}
			if err := record.Unmarshal(v); err != nil {
				return err
			}

			result = append(result, record)
		}

		return nil
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
	var m PersistedEvent
	m.Id = e.ID[:]
	m.Stream = e.Stream[:]
	m.Version = e.Version
	m.Type = e.Type
	m.Data = e.Data
	m.Metadata = e.Metadata
	m.CausationId = e.CausationID[:]
	m.CorrelationId = e.CorrelationID[:]
	m.AddedAt = e.AddedAt.UnixNano()

	return proto.Marshal(&m)
}

func (e *Event) Unmarshal(data []byte) error {
	var m PersistedEvent
	if err := proto.Unmarshal(data, &m); err != nil {
		return err
	}

	var id ulid.ULID
	if err := id.UnmarshalBinary(m.Id); err != nil {
		return err
	}
	e.ID = id

	var stream uuid.UUID
	if err := stream.UnmarshalBinary(m.Stream); err != nil {
		return err
	}
	e.Stream = stream

	e.Version = m.Version
	e.Type = m.Type
	e.Data = m.Data
	e.Metadata = m.Metadata

	var causationID ulid.ULID
	if err := causationID.UnmarshalBinary(m.CausationId); err != nil {
		return err
	}
	e.CausationID = causationID

	var correlationID ulid.ULID
	if err := correlationID.UnmarshalBinary(m.CorrelationId); err != nil {
		return err
	}
	e.CorrelationID = correlationID

	e.AddedAt = time.Unix(0, m.AddedAt)

	return nil
}

type Stream struct {
	ID      uuid.UUID
	Events  []ulid.ULID
	AddedAt time.Time
}

func (s *Stream) Marshal() ([]byte, error) {
	var m PersistedStream

	m.Id = s.ID[:]

	for _, event := range s.Events {
		id, err := event.MarshalBinary()
		if err != nil {
			return nil, err
		}
		m.Events = append(m.Events, id)
	}

	m.AddedAt = s.AddedAt.UnixNano()

	return proto.Marshal(&m)
}

func (s *Stream) Unmarshal(data []byte) error {
	var m PersistedStream

	if err := proto.Unmarshal(data, &m); err != nil {
		return err
	}

	var id uuid.UUID
	if err := id.UnmarshalBinary(m.Id); err != nil {
		return err
	}

	s.ID = id
	s.Events = []ulid.ULID{}

	for _, id := range m.Events {
		var event ulid.ULID
		if err := event.UnmarshalBinary(id); err != nil {
			return err
		}
		s.Events = append(s.Events, event)
	}

	s.AddedAt = time.Unix(0, m.AddedAt)

	return nil
}
