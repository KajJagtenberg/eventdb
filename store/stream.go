package store

import (
	"bytes"

	"github.com/oklog/ulid"
)

type Stream struct {
	Events []ulid.ULID
}

func (stream *Stream) Marshal() []byte {
	buffer := new(bytes.Buffer)

	for _, id := range stream.Events {
		buffer.Write(id[:])
	}

	return buffer.Bytes()
}

func (stream *Stream) Unmarshal(data []byte) {
	stream.Events = []ulid.ULID{}

	for len(data) >= 16 {
		var event ulid.ULID

		copy(event[:], data[:16])
		data = data[16:]

		stream.Events = append(stream.Events, event)
	}
}
