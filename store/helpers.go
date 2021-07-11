package store

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/eventflowdb/eventflowdb/api"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func getStream(txn *badger.Txn, stream uuid.UUID) (*PersistedStream, error) {
	var persistedStream PersistedStream

	item, err := txn.Get(getStreamKey(stream))
	switch err {
	case nil:
		var value Value
		if err := item.Value(func(val []byte) error {
			return proto.Unmarshal(val, &value)
		}); err != nil {
			return nil, err
		}

		if err := proto.Unmarshal(value.Data, &persistedStream); err != nil {
			return nil, err
		}
	case badger.ErrKeyNotFound:
		persistedStream.Id = stream.String()
		persistedStream.Version = 0
		persistedStream.CreatedAt = time.Now().Unix()
	default:
		return nil, err
	}

	return &persistedStream, nil
}

func setStream(txn *badger.Txn, id uuid.UUID, stream *PersistedStream) error {
	data, err := proto.Marshal(stream)
	if err != nil {
		return err
	}
	value, err := proto.Marshal(&Value{
		Version:  0,
		Encoding: 1,
		Data:     data,
	})
	if err != nil {
		return err
	}
	return txn.Set(getStreamKey(id), value)
}

func getStreamEventKey(stream uuid.UUID, version uint32) []byte {
	buf := new(bytes.Buffer)
	buf.Write(PREFIX_STREAM_EVENT)
	buf.Write(SEPERATOR)
	buf.WriteString(stream.String())
	buf.Write(SEPERATOR)

	binary.Write(buf, binary.BigEndian, version)

	return buf.Bytes()
}

func getStreamKey(stream uuid.UUID) []byte {
	buf := new(bytes.Buffer)
	buf.Write(PREFIX_STREAM_METADATA)
	buf.Write(SEPERATOR)
	buf.WriteString(stream.String())

	return buf.Bytes()
}

func getEventKey(id uuid.UUID) []byte {
	buf := new(bytes.Buffer)
	buf.Write(PREFIX_EVENT)
	buf.Write(SEPERATOR)
	buf.WriteString(id.String())
	buf.Write(SEPERATOR)

	return buf.Bytes()
}

func setEvent(txn *badger.Txn, id uuid.UUID, stream uuid.UUID, event *api.EventData, version uint32) error {
	data, err := proto.Marshal(&api.Event{
		Id:            id.String(),
		Stream:        stream.String(),
		Version:       version,
		Type:          event.Type,
		Metadata:      event.Metadata,
		CausationId:   event.CausationId,
		CorrelationId: event.CorrelationId,
		Data:          event.Data,
		AddedAt:       time.Now().Unix(),
	})
	if err != nil {
		return err
	}
	value, err := proto.Marshal(&Value{
		Encoding: 1,
		Version:  0,
		Data:     data,
	})
	if err != nil {
		return err
	}

	return txn.Set(getEventKey(id), value)
}

func setStreamEvent(txn *badger.Txn, stream uuid.UUID, version uint32, event uuid.UUID) error {
	data, err := proto.Marshal(&StreamEvent{
		Id: event.String(),
	})
	if err != nil {
		return err
	}
	value, err := proto.Marshal(&Value{
		Version:  0,
		Encoding: 1,
		Data:     data,
	})
	if err != nil {
		return err
	}

	return txn.Set(getStreamEventKey(stream, version), value)
}

func getStreamEvent(txn *badger.Txn, stream uuid.UUID, version uint32) (string, error) {
	item, err := txn.Get(getStreamEventKey(stream, version))
	if err != nil {
		return "", err
	}

	var id string

	if err := item.Value(func(val []byte) error {
		var value Value
		if err := proto.Unmarshal(val, &value); err != nil {
			return err
		}
		var event StreamEvent
		if err := proto.Unmarshal(value.Data, &event); err != nil {
			return err
		}

		id = event.Id

		return nil
	}); err != nil {
		return "", nil
	}

	return id, nil
}
