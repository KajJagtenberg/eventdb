package store

import (
	protos "eventflowdb/gen"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

type Stream struct {
	ID        uuid.UUID   `json:"name"`
	Events    []ulid.ULID `json:"events"`
	CreatedAt time.Time   `json:"created_at"`
}

func (stream *Stream) Size() int {
	return len(stream.Events)
}

func (stream *Stream) Serialize() ([]byte, error) {
	t, err := stream.CreatedAt.MarshalBinary()
	if err != nil {
		return nil, err
	}

	var events [][]byte

	for _, event := range stream.Events {
		events = append(events, event[:])
	}

	m := protos.Stream{
		Id:      stream.ID[:],
		AddedAt: t,
		Events:  events,
	}

	return proto.Marshal(&m)
}

func (stream *Stream) Deserialize(data []byte) error {
	var m protos.Stream

	if err := proto.Unmarshal(data, &m); err != nil {
		return err
	}

	id, err := uuid.FromBytes(m.Id)
	if err != nil {
		return err
	}

	stream.ID = id

	if err := stream.CreatedAt.UnmarshalBinary(m.AddedAt); err != nil {
		return err
	}

	for _, event := range m.Events {
		var id ulid.ULID
		if err := id.UnmarshalBinary(event); err != nil {
			return err
		}

		stream.Events = append(stream.Events, id)
	}

	return nil
}
