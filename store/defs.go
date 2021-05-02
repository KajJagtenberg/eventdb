package store

import (
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

type EventData struct {
	Type          string    `json:"type"`
	Data          []byte    `json:"data"`
	Metadata      []byte    `json:"metadata"`
	CausationID   ulid.ULID `json:"causation_id"`
	CorrelationID ulid.ULID `json:"correlation_id"`
	// AddedAt       time.Time `json:"added_at"`
}

type Event struct {
	ID            ulid.ULID `json:"id"`
	Stream        uuid.UUID `json:"stream"`
	Version       uint32    `json:"version"`
	Type          string    `json:"type"`
	Data          []byte    `json:"data"`
	Metadata      []byte    `json:"metadata"`
	CausationID   ulid.ULID `json:"causation_id"`
	CorrelationID ulid.ULID `json:"correlation_id"`
	AddedAt       time.Time `json:"added_at"`
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
