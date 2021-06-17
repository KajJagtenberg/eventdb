package store

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/conv"
	"github.com/oklog/ulid"
	"google.golang.org/protobuf/proto"
)

type BadgerEventStore struct {
	db                  *badger.DB
	estimateStreamCount int64
	estimateEventCount  int64
	lock                sync.Mutex
}

var (
	BUCKET_EVENTS   = []byte{0, 0}
	BUCKET_STREAMS  = []byte{0, 1}
	BUCKET_METADATA = []byte{0, 2}
)

func (s *BadgerEventStore) Add(req *api.AddRequest) (res *api.EventResponse, err error) {
	res = &api.EventResponse{
		Events: make([]*api.Event, 0),
	}

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

	if bytes.Equal(stream[:], make([]byte, 16)) {
		return nil, errors.New("stream cannot be all zeroes")
	}

	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	var persistedStream PersistedStream

	item, err := txn.Get(getStreamKey(stream))
	if err == nil {
		data, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}

		err = proto.Unmarshal(data, &persistedStream)
		if err != nil {
			return nil, err
		}
	} else if err == badger.ErrKeyNotFound {
		persistedStream.Id = stream[:]
		persistedStream.AddedAt = time.Now().Unix()
	} else {
		return nil, err
	}

	if int(req.Version) < len(persistedStream.Events) {
		return nil, ErrConcurrentStreamModification
	}

	if int(req.Version) > len(persistedStream.Events) {
		return nil, ErrGappedStream
	}

	now := time.Now()

	s.lock.Lock()

	for i, event := range req.Events {
		if len(event.Type) == 0 {
			return nil, ErrEmptyEventType
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

		if err := txn.Set(getEventKey(id), data); err != nil {
			return nil, err
		}

		res.Events = append(res.Events, &api.Event{
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

	s.lock.Unlock()

	data, err := proto.Marshal(&persistedStream)
	if err != nil {
		return nil, err
	}

	if err := txn.Set(getStreamKey(stream), data); err != nil {
		return nil, err
	}

	return res, txn.Commit()
}

func (s *BadgerEventStore) Get(req *api.GetRequest) (res *api.EventResponse, err error) {
	res = &api.EventResponse{
		Events: make([]*api.Event, 0),
	}

	stream, err := uuid.Parse(req.Stream)
	if err != nil {
		return nil, err
	}

	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get(getStreamKey(stream))

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return res, nil
		} else {
			return nil, err
		}
	}

	data, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	var persistedStream PersistedStream
	if err := proto.Unmarshal(data, &persistedStream); err != nil {
		return nil, err
	}

	for _, key := range persistedStream.Events {
		if req.Version > 0 {
			req.Version--
			continue
		}

		if len(res.Events) >= int(req.Limit) && req.Limit != 0 {
			break
		}

		var id ulid.ULID
		if err := id.UnmarshalBinary(key); err != nil {
			return nil, err
		}

		item, err := txn.Get(getEventKey(id))
		if err != nil {
			return nil, err
		}

		data, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}

		var event PersistedEvent
		if err := proto.Unmarshal(data, &event); err != nil {
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

		res.Events = append(res.Events, &api.Event{
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

	return res, nil
}

func (s *BadgerEventStore) GetAll(req *api.GetAllRequest) (res *api.EventResponse, err error) {
	res = &api.EventResponse{
		Events: make([]*api.Event, 0),
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	var offset ulid.ULID

	if len(req.Offset) > 0 {
		offset, err = ulid.Parse(req.Offset)
		if err != nil {
			return nil, err
		}
	}

	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	cursor := txn.NewIterator(badger.DefaultIteratorOptions)
	defer cursor.Close()

	for cursor.Seek(getEventKey(offset)); cursor.ValidForPrefix(BUCKET_EVENTS); cursor.Next() {
		item := cursor.Item()

		if bytes.Equal(item.Key(), offset[:]) {
			continue
		}

		if len(res.Events) >= int(req.Limit) {
			break
		}

		value, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}

		var event PersistedEvent
		if err := proto.Unmarshal(value, &event); err != nil {
			return nil, err
		}

		var id ulid.ULID
		if err := id.UnmarshalBinary(item.Key()[2:]); err != nil {
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

		res.Events = append(res.Events, &api.Event{
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

	return res, nil
}

func (s *BadgerEventStore) EventCount(req *api.EventCountRequest) (res *api.EventCountResponse, err error) {
	res = &api.EventCountResponse{}

	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	cursor := txn.NewIterator(badger.DefaultIteratorOptions)

	defer cursor.Close()

	for cursor.Seek(BUCKET_EVENTS); cursor.ValidForPrefix(BUCKET_EVENTS); cursor.Next() {
		res.Count++
	}

	return res, nil
}

func (s *BadgerEventStore) StreamCount(req *api.StreamCountRequest) (res *api.StreamCountResponse, err error) {
	res = &api.StreamCountResponse{}

	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	cursor := txn.NewIterator(badger.DefaultIteratorOptions)

	defer cursor.Close()

	for cursor.Seek(BUCKET_STREAMS); cursor.ValidForPrefix(BUCKET_STREAMS); cursor.Next() {
		res.Count++
	}

	return res, nil
}

func (s *BadgerEventStore) EventCountEstimate(req *api.EventCountEstimateRequest) (res *api.EventCountResponse, err error) {
	return &api.EventCountResponse{
		Count: s.estimateEventCount,
	}, err
}

func (s *BadgerEventStore) StreamCountEstimate(req *api.StreamCountEstimateRequest) (res *api.StreamCountResponse, err error) {
	return &api.StreamCountResponse{
		Count: s.estimateStreamCount,
	}, err
}

func (s *BadgerEventStore) Size(req *api.SizeRequest) (res *api.SizeResponse, err error) {
	res = &api.SizeResponse{}

	lsm, _ := s.db.Size()

	res.Size = lsm
	res.SizeHuman = conv.ByteCountSI(res.Size)

	return res, nil
}

func (s *BadgerEventStore) ListStreams(req *api.ListStreamsRequest) (res *api.ListStreamsReponse, err error) {
	res = &api.ListStreamsReponse{}

	if req.Limit == 0 {
		req.Limit = 10
	}

	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	prefix := BUCKET_STREAMS

	cursor := txn.NewIterator(badger.DefaultIteratorOptions)
	defer cursor.Close()

	for cursor.Seek(prefix); cursor.ValidForPrefix(prefix); cursor.Next() {
		if req.Skip > 0 {
			req.Skip--
			continue
		}

		if len(res.Streams) >= int(req.Limit) {
			return res, nil
		}

		value, err := cursor.Item().ValueCopy(nil)
		if err != nil {
			return nil, err
		}

		var persistedStream PersistedStream
		if err := proto.Unmarshal(value, &persistedStream); err != nil {
			return nil, err
		}

		var id uuid.UUID
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

	return res, nil
}

func (s *BadgerEventStore) Backup(dst io.Writer) error {
	_, err := s.db.Backup(dst, 0)
	return err
}

func (s *BadgerEventStore) Close() error {
	return s.db.Close()
}

type BadgerStoreOptions struct {
	DB             *badger.DB
	EstimateCounts bool
}

func NewBadgerEventStore(options BadgerStoreOptions) (*BadgerEventStore, error) {
	db := options.DB

	store := &BadgerEventStore{db, 0, 0, sync.Mutex{}}

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

	if !db.Opts().InMemory {
		go func() {
			if err := db.RunValueLogGC(0.7); err != nil && err != badger.ErrNoRewrite {
				log.Fatal(err)
			}

			time.Sleep(ESTIMATE_SLEEP_TIME)
		}()
	}

	return store, nil
}

func getStreamKey(stream uuid.UUID) []byte {
	result := BUCKET_STREAMS
	result = append(result, stream[:]...)
	return result
}

func getEventKey(id ulid.ULID) []byte {
	result := BUCKET_EVENTS
	result = append(result, id[:]...)
	return result
}
