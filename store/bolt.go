package store

import (
	"bytes"
	"crypto/rand"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/oklog/ulid"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

type boltEventStore struct {
	db *bbolt.DB
}

func (s *boltEventStore) Add(req *api.AddRequest) (*api.EventResponse, error) {
	stream, err := uuid.Parse(req.Stream)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(stream[:], make([]byte, 16)) {
		return nil, errors.New("stream cannot be all zeroes")
	}

	if len(req.Events) == 0 {
		return nil, errors.New("list of events is empty")
	}

	res := &api.EventResponse{}

	txn, err := s.db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	streams := txn.Bucket([]byte("streams"))
	events := txn.Bucket([]byte("events"))

	var persistedStream PersistedStream

	if data := streams.Get(stream[:]); err != nil {
		if err := proto.Unmarshal(data, &persistedStream); err != nil {
			return nil, err
		}
	} else {
		persistedStream.Id = stream[:]
		persistedStream.AddedAt = time.Now().Unix()
	}

	if int(req.Version) < len(persistedStream.Events) {
		return nil, ErrConcurrentStreamModification
	}

	if int(req.Version) > len(persistedStream.Events) {
		return nil, ErrGappedStream
	}

	now := time.Now()

	for i, event := range req.Events {
		if len(event.Type) == 0 {
			return nil, errors.New("event type cannot be empty")
		}

		id, err := ulid.New(ulid.Now(), rand.Reader)
		if err != nil {
			return nil, err
		}

		causationId, err := ulid.Parse(event.CausationId)
		if err != nil {
			return nil, err
		}

		correlationId, err := ulid.Parse(event.CorrelationId)
		if err != nil {
			return nil, err
		}

		record := PersistedEvent{
			Id:            id[:],
			Stream:        stream[:],
			Version:       req.Version + uint32(i),
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   causationId[:],
			CorrelationId: correlationId[:],
			AddedAt:       now.Unix(),
		}

		data, err := proto.Marshal(&record)
		if err != nil {
			return nil, err
		}

		if err := events.Put(id[:], data); err != nil {
			return nil, err
		}

		res.Events = append(res.Events, &api.EventResponse_Event{
			Id:            id.String(),
			Stream:        stream.String(),
			Version:       req.Version + uint32(i),
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   event.CorrelationId,
			CorrelationId: event.CorrelationId,
			AddedAt:       now.Unix(),
		})

		persistedStream.Events = append(persistedStream.Events, record.Id)
	}

	data, err := proto.Marshal(&persistedStream)
	if err != nil {
		return nil, err
	}

	if err := streams.Put(stream[:], data); err != nil {
		return nil, err
	}

	return res, txn.Commit()
}

func (s *boltEventStore) Get(req *api.GetRequest) (res *api.EventResponse, err error) {
	stream, err := uuid.Parse(req.Stream)
	if err != nil {
		return nil, err
	}

	txn, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	streams := txn.Bucket([]byte("streams"))
	events := txn.Bucket([]byte("streams"))

	data := streams.Get(stream[:])
	if data == nil {
		return res, nil
	}

	var persistedStream PersistedStream
	if err := proto.Unmarshal(data, &persistedStream); err != nil {
		return nil, err
	}

	for _, key := range persistedStream.Events {
		data := events.Get(key)

		var event PersistedEvent
		if err := proto.Unmarshal(data, &event); err != nil {
			return nil, err
		}

		var id ulid.ULID
		if err := id.UnmarshalBinary(key); err != nil {
			return nil, err
		}

		var stream uuid.UUID
		if err := stream.UnmarshalBinary(event.Stream); err != nil {
			return nil, err
		}

		var causationId ulid.ULID
		if err := causationId.UnmarshalBinary(event.CausationId); err != nil {
			return nil, err
		}

		var correlationId ulid.ULID
		if err := correlationId.UnmarshalBinary(event.CorrelationId); err != nil {
			return nil, err
		}

		res.Events = append(res.Events, &api.EventResponse_Event{
			Id:            id.String(),
			Stream:        stream.String(),
			Version:       event.Version,
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   causationId.String(),
			CorrelationId: correlationId.String(),
			AddedAt:       event.AddedAt,
		})
	}

	return res, txn.Commit()
}

func (s *boltEventStore) GetAll(req *api.GetAllRequest) (res *api.EventResponse, err error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	offset, err := ulid.Parse(req.Offset)
	if err != nil {
		return nil, err
	}

	txn, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	events := txn.Bucket([]byte("streams"))

	cursor := events.Cursor()

	for k, v := cursor.Seek(offset[:]); k != nil; k, v = cursor.Next() {
		if bytes.Equal(k, offset[:]) {
			continue
		}

		if len(res.Events) >= int(req.Limit) {
			break
		}

		var event PersistedEvent
		if err := proto.Unmarshal(v, &event); err != nil {
			return nil, err
		}

		var id ulid.ULID
		if err := id.UnmarshalBinary(k); err != nil {
			return nil, err
		}

		var stream uuid.UUID
		if err := stream.UnmarshalBinary(event.Stream); err != nil {
			return nil, err
		}

		var causationId ulid.ULID
		if err := causationId.UnmarshalBinary(event.CausationId); err != nil {
			return nil, err
		}

		var correlationId ulid.ULID
		if err := correlationId.UnmarshalBinary(event.CorrelationId); err != nil {
			return nil, err
		}

		res.Events = append(res.Events, &api.EventResponse_Event{
			Id:            id.String(),
			Stream:        stream.String(),
			Version:       event.Version,
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   causationId.String(),
			CorrelationId: correlationId.String(),
			AddedAt:       event.AddedAt,
		})
	}

	return res, txn.Commit()
}

func NewBoltEventStore(db *bbolt.DB) (*boltEventStore, error) {
	txn, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	for _, bucket := range []string{"events", "streams"} {
		if _, err := txn.CreateBucketIfNotExists([]byte(bucket)); err != nil {
			return nil, err
		}
	}

	return &boltEventStore{db}, txn.Commit()
}
