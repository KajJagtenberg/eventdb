package store

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/si"
	"github.com/oklog/ulid"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

const (
	ESTIMATE_SLEEP_TIME = time.Second
)

var (
	ErrConcurrentStreamModification = errors.New("concurrent stream modification")
	ErrGappedStream                 = errors.New("given version leaves gap in stream")
)

type boltEventStore struct {
	db                  *bbolt.DB
	estimateStreamCount int64
	estimateEventCount  int64
}

func (s *boltEventStore) Add(req *api.AddRequest) (res *api.EventResponse, err error) {
	res = &api.EventResponse{}

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

		var causationId ulid.ULID
		var correlationId ulid.ULID

		if len(event.CausationId) > 0 {
			causationId, err = ulid.Parse(event.CausationId)
			if err != nil {
				return nil, err
			}
		}

		if len(event.CorrelationId) > 0 {
			correlationId, err = ulid.Parse(event.CorrelationId)
			if err != nil {
				return nil, err
			}
		}

		if causationId.String() == "00000000000000000000000000" {
			causationId = id
		}

		if correlationId.String() == "00000000000000000000000000" {
			correlationId = id
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
			CausationId:   causationId.String(),
			CorrelationId: correlationId.String(),
			AddedAt:       now.Unix(),
		})

		persistedStream.Events = append(persistedStream.Events, id[:])
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
	res = &api.EventResponse{}

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
	events := txn.Bucket([]byte("events"))

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

	return res, txn.Rollback()
}

func (s *boltEventStore) GetAll(req *api.GetAllRequest) (res *api.EventResponse, err error) {
	res = &api.EventResponse{}

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

	return res, txn.Rollback()
}

func (s *boltEventStore) EventCount(req *api.EventCountRequest) (res *api.EventCountResponse, err error) {
	res = &api.EventCountResponse{}

	txn, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	cursor := txn.Bucket([]byte("events")).Cursor()

	for k, _ := cursor.First(); k != nil; cursor.Next() {
		res.Count++
	}

	return res, txn.Commit()
}

func (s *boltEventStore) StreamCount(req *api.StreamCountRequest) (res *api.StreamCountResponse, err error) {
	res = &api.StreamCountResponse{}

	txn, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	cursor := txn.Bucket([]byte("streams")).Cursor()

	for k, _ := cursor.First(); k != nil; cursor.Next() {
		res.Count++
	}

	return res, txn.Rollback()
}

func (s *boltEventStore) EventCountEstimate(req *api.EventCountEstimateRequest) (res *api.EventCountResponse, err error) {
	return &api.EventCountResponse{
		Count: s.estimateEventCount,
	}, err
}

func (s *boltEventStore) StreamCountEstimate(req *api.StreamCountEstimateRequest) (res *api.StreamCountResponse, err error) {
	return &api.StreamCountResponse{
		Count: s.estimateStreamCount,
	}, err
}

func (s *boltEventStore) Size(req *api.SizeRequest) (res *api.SizeResponse, err error) {
	res = &api.SizeResponse{}

	txn, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	res.Size = txn.Size()
	res.SizeHuman = si.ByteCountSI(res.Size)

	return res, txn.Rollback()
}

func (s *boltEventStore) ListStreams(req *api.ListStreamsRequest) (res *api.ListStreamsReponse, err error) {
	res = &api.ListStreamsReponse{}

	if req.Limit == 0 {
		req.Limit = 10
	}

	txn, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	cursor := txn.Bucket([]byte("streams")).Cursor()

	for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
		if req.Skip > 0 {
			req.Skip--
			continue
		}

		if len(res.Streams) >= int(req.Limit) {
			break
		}

		var persistedStream PersistedStream
		if err := proto.Unmarshal(v, &persistedStream); err != nil {
			return nil, err
		}

		var id ulid.ULID
		if err := id.UnmarshalBinary(persistedStream.Id); err != nil {
			return nil, err
		}

		var events []string

		for _, idByte := range persistedStream.Events {
			var id ulid.ULID
			if err := id.UnmarshalBinary(idByte); err != nil {
				return nil, err
			}

			events = append(events, id.String())
		}

		res.Streams = append(res.Streams, &api.ListStreamsReponse_Stream{
			Id:      id.String(),
			Events:  events,
			AddedAt: persistedStream.AddedAt,
		})
	}

	return res, txn.Rollback()
}

func (s *boltEventStore) Backup(dst io.Writer) error {
	txn, err := s.db.Begin(false)
	if err != nil {
		return err
	}
	defer txn.Rollback()

	if _, err := txn.WriteTo(dst); err != nil {
		return err
	}

	return txn.Rollback()
}

func (s *boltEventStore) Close() error {
	return s.db.Close()
}

func NewBoltEventStore(options BoltStoreOptions) (*boltEventStore, error) {
	db := options.DB

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

	if err := txn.Commit(); err != nil {
		return nil, err
	}

	store := &boltEventStore{db, 0, 0}

	if options.EstimateCounts {
		go func() {
			for {
				streamCount, err := store.StreamCount(&api.StreamCountRequest{})
				if err != nil {
					log.Fatalf("failed to get stream count: %v", err)
				}

				eventCount, err := store.EventCount(&api.EventCountRequest{})
				if err != nil {
					log.Fatalf("failed to get stream count: %v", err)
				}

				store.estimateStreamCount = streamCount.Count
				store.estimateEventCount = eventCount.Count

				time.Sleep(ESTIMATE_SLEEP_TIME)
			}
		}()
	}

	return store, err
}

type BoltStoreOptions struct {
	DB             *bbolt.DB
	EstimateCounts bool
}
