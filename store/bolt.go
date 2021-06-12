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
