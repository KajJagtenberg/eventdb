package store

import (
	protos "eventflowdb/gen"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

type EventData struct {
	Type          string    `json:"type"`
	Data          []byte    `json:"data"`
	Metadata      []byte    `json:"metadata"`
	CausationID   ulid.ULID `json:"causation_id"`
	CorrelationID ulid.ULID `json:"correlation_id"`
}

type RecordedEvent struct {
	ID            ulid.ULID `json:"id"`
	Stream        uuid.UUID `json:"stream"`
	Version       uint32    `json:"version"`
	Type          string    `json:"type"`
	Data          []byte    `json:"data"`
	Metadata      []byte    `json:"metadata"`
	AddedAt       time.Time `json:"added_at"`
	CausationID   ulid.ULID `json:"causation_id"`
	CorrelationID ulid.ULID `json:"correlation_id"`
}

func (event *RecordedEvent) Serialize() ([]byte, error) {
	t, err := event.AddedAt.MarshalBinary()
	if err != nil {
		return nil, err
	}

	m := protos.RecordedEvent{
		Id:            event.ID[:],
		Stream:        event.Stream[:],
		Version:       event.Version,
		Type:          event.Type,
		Data:          event.Data,
		Metadata:      event.Metadata,
		AddedAt:       t,
		CausationId:   event.CausationID[:],
		CorrelationId: event.CorrelationID[:],
	}

	return proto.Marshal(&m)
}

func (event *RecordedEvent) Deserialize(data []byte) error {
	m := protos.RecordedEvent{}

	if err := proto.Unmarshal(data, &m); err != nil {
		return err
	}

	if err := event.ID.UnmarshalBinary(m.Id); err != nil {
		return err
	}

	stream, err := uuid.FromBytes(m.Stream)
	if err != nil {
		return err
	}

	if err := event.CausationID.UnmarshalBinary(m.CausationId); err != nil {
		return err
	}

	if err := event.CorrelationID.UnmarshalBinary(m.CorrelationId); err != nil {
		return err
	}

	if err := event.AddedAt.UnmarshalBinary(m.AddedAt); err != nil {
		return err
	}

	event.Stream = stream
	event.Version = m.Version
	event.Type = m.Type
	event.Data = m.Data
	event.Metadata = m.Metadata

	return nil
}
