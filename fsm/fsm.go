package fsm

import (
	"io"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/kajjagtenberg/eventflowdb/transport"
	"github.com/oklog/ulid"
	"google.golang.org/protobuf/proto"
)

type badgerFSM struct {
	db    *badger.DB
	store store.EventStore
}

func (b *badgerFSM) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:
		var cmd = Command{}
		if err := proto.Unmarshal(log.Data, &cmd); err != nil {
			return nil
		}

		op := strings.ToUpper(strings.TrimSpace(cmd.Op))

		switch op {
		case "ADD":
			data, err := b.add(cmd.Payload)
			return &ApplyResponse{
				Error: err,
				Data:  data,
			}
		}
	}

	return nil
}

func (b *badgerFSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (b *badgerFSM) Restore(io.ReadCloser) error {
	return nil
}
func (b *badgerFSM) add(cmd []byte) (interface{}, error) {
	var req transport.AddRequest
	if err := proto.Unmarshal(cmd, &req); err != nil {
		return nil, err
	}

	stream, err := uuid.Parse(req.Stream)
	if err != nil {
		return nil, err
	}

	var eventdata []store.EventData

	for _, e := range req.Events {
		causation_id, err := ulid.Parse(e.CausationId)
		if err != nil {
			return nil, err
		}
		correlation_id, err := ulid.Parse(e.CausationId)
		if err != nil {
			return nil, err
		}

		eventdata = append(eventdata, store.EventData{
			Type:          e.Type,
			Data:          e.Data,
			Metadata:      e.Metadata,
			CausationID:   causation_id,
			CorrelationID: correlation_id,
		})
	}

	events, err := b.store.Add(stream, req.Version, eventdata)
	if err != nil {
		return nil, err
	}

	var parsedEvent []*transport.EventResponse_Event

	for _, event := range events {
		parsedEvent = append(parsedEvent, &transport.EventResponse_Event{
			Id:            event.ID.String(),
			Stream:        event.Stream.String(),
			Version:       event.Version,
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   event.CausationID.String(),
			CorrelationId: event.CorrelationID.String(),
			AddedAt:       event.AddedAt.Unix(),
		})
	}

	return &transport.EventResponse{
		Events: parsedEvent,
	}, nil
}

func NewBadgerFSM(db *badger.DB, store store.EventStore) *badgerFSM {
	return &badgerFSM{db, store}
}
