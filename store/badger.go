package store

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/eventflowdb/eventflowdb/api"
	"github.com/eventflowdb/eventflowdb/conv"
	"github.com/google/uuid"
	"github.com/karlseguin/ccache/v2"
	"github.com/oklog/ulid"
	"google.golang.org/protobuf/proto"
)

type BadgerEventStore struct {
	db *badger.DB
	// estimateStreamCount int64
	// estimateEventCount  int64
	cache *ccache.Cache
}

var (
	BUCKET_EVENTS      = []byte{0}
	BUCKET_STREAMS     = []byte{1}
	BUCKET_SEQUENCE    = []byte{2}
	BUCKET_SYSTEM      = []byte{3}
	BUCKET_STREAM_LIST = []byte{4}

	KEY_CURRENT_SEQUENCE = []byte{3, 0}
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

	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	if err := txn.Set(getStreamListKey(stream), nil); err != nil {
		return nil, err
	}

	// Retrieve the current global event sequence
	sequence, err := getCurrentSequence(txn)
	if err != nil {
		return nil, err
	}

	// Check if the gap does not exist
	if req.Version > 0 {
		_, err := txn.Get(getStreamVersionKey(stream, req.Version-1))
		if err == badger.ErrKeyNotFound {
			return nil, ErrGappedStream
		}
	}

	for i, event := range req.Events {
		// Validate event
		if len(event.Type) == 0 {
			return nil, ErrEmptyEventType
		}

		var id ulid.ULID

		// Check if the given event id already exists, otherwise generate it
		if len(event.Id) > 0 {
			id, err = ulid.Parse(event.Id)
			if err != nil {
				return nil, err
			}

			_, err := txn.Get(getEventKey(id))
			if err == badger.ErrKeyNotFound {
				continue
			}
		} else {
			id, err = ulid.New(ulid.Now(), rand.Reader)
			if err != nil {
				return nil, err
			}
		}

		version := req.Version + uint32(i)

		// Check if [stream][version] is empty
		_, err := txn.Get(getStreamVersionKey(stream, version))
		if err != badger.ErrKeyNotFound {
			return nil, ErrConcurrentStreamModification
		}

		// Get causation id from the given event, otherwise assign the event id to it
		var causationID ulid.ULID
		if len(event.CausationId) > 0 {
			causationID, err = ulid.Parse(event.CausationId)
			if err != nil {
				return nil, err
			}
		} else {
			causationID = id
		}

		// Get correlation id from the given event, otherwise assign the event id to it
		var correlationID ulid.ULID
		if len(event.CorrelationId) > 0 {
			correlationID, err = ulid.Parse(event.CorrelationId)
			if err != nil {
				return nil, err
			}
		} else {
			correlationID = id
		}

		// Marshal the event
		data, err := proto.Marshal(&PersistedEvent{
			Id:            id[:],
			Stream:        stream[:],
			Version:       version,
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   causationID[:],
			CorrelationId: correlationID[:],
		})
		if err != nil {
			return nil, err
		}

		//Persist the event
		if err := txn.Set(getEventKey(id), data); err != nil {
			return nil, err
		}

		// Persist the event id to the [stream][version]
		if err := txn.Set(getStreamVersionKey(stream, version), id[:]); err != nil {
			return nil, err
		}

		// Persist the event id to the [sequence]
		if err := txn.Set(getSequenceKey(sequence), id[:]); err != nil {
			return nil, err
		}

		sequence++

		res.Events = append(res.Events, &api.Event{
			Id:            id.String(),
			Stream:        stream.String(),
			Version:       version,
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   causationID.String(),
			CorrelationId: correlationID.String(),
		})
	}

	setCurrentSequence(txn, sequence)

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

	cursor := txn.NewIterator(badger.DefaultIteratorOptions)
	defer cursor.Close()

	for cursor.Seek(getStreamVersionKey(stream, req.Version)); cursor.ValidForPrefix(getStreamKey(stream)); cursor.Next() {
		if len(res.Events) >= int(req.Limit) && req.Limit != 0 {
			break
		}

		if err := cursor.Item().Value(func(val []byte) error {
			var id ulid.ULID
			if err := id.UnmarshalBinary(val); err != nil {
				return err
			}

			item, err := txn.Get(getEventKey(id))
			if err != nil {
				return err
			}

			return item.Value(func(val []byte) error {
				var event PersistedEvent
				if err := proto.Unmarshal(val, &event); err != nil {
					return err
				}

				var id ulid.ULID
				if err := id.UnmarshalBinary(event.Id); err != nil {
					return err
				}

				var causationID ulid.ULID
				if err := causationID.UnmarshalBinary(event.CausationId); err != nil {
					return err
				}

				var correlationID ulid.ULID
				if err := correlationID.UnmarshalBinary(event.CorrelationId); err != nil {
					return err
				}

				res.Events = append(res.Events, &api.Event{
					Id:            id.String(),
					Stream:        stream.String(),
					Version:       event.Version,
					Type:          event.Type,
					Data:          event.Data,
					Metadata:      event.Metadata,
					CausationId:   causationID.String(),
					CorrelationId: correlationID.String(),
				})

				return nil
			})
		}); err != nil {
			return nil, err
		}
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

	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	cursor := txn.NewIterator(badger.DefaultIteratorOptions)
	defer cursor.Close()

	for cursor.Seek(getSequenceKey(req.Offset)); cursor.ValidForPrefix(BUCKET_SEQUENCE); cursor.Next() {
		if len(res.Events) >= int(req.Limit) {
			break
		}

		if err := cursor.Item().Value(func(val []byte) error {
			var id ulid.ULID
			if err := id.UnmarshalBinary(val); err != nil {
				return err
			}

			item, err := txn.Get(getEventKey(id))
			if err != nil {
				return err
			}

			return item.Value(func(val []byte) error {
				var event PersistedEvent
				if err := proto.Unmarshal(val, &event); err != nil {
					return err
				}

				var stream uuid.UUID
				if err := stream.UnmarshalBinary(event.Stream); err != nil {
					return err
				}

				var id ulid.ULID
				if err := id.UnmarshalBinary(event.Id); err != nil {
					return err
				}

				var causationID ulid.ULID
				if err := causationID.UnmarshalBinary(event.CausationId); err != nil {
					return err
				}

				var correlationID ulid.ULID
				if err := correlationID.UnmarshalBinary(event.CorrelationId); err != nil {
					return err
				}

				res.Events = append(res.Events, &api.Event{
					Id:            id.String(),
					Stream:        stream.String(),
					Version:       event.Version,
					Type:          event.Type,
					Data:          event.Data,
					Metadata:      event.Metadata,
					CausationId:   causationID.String(),
					CorrelationId: correlationID.String(),
				})

				return nil
			})
		}); err != nil {
			return nil, err
		}
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
	item, err := s.cache.Fetch("user:4", ESTIMATE_TTL, func() (interface{}, error) {
		res, err := s.EventCount(&api.EventCountRequest{})
		if err != nil {
			return nil, err
		}
		return res.Count, nil
	})
	if err != nil {
		return nil, err
	}

	return &api.EventCountResponse{
		Count: item.Value().(int64),
	}, err
}

func (s *BadgerEventStore) StreamCountEstimate(req *api.StreamCountEstimateRequest) (res *api.StreamCountResponse, err error) {
	item, err := s.cache.Fetch("user:4", ESTIMATE_TTL, func() (interface{}, error) {
		res, err := s.StreamCount(&api.StreamCountRequest{})
		if err != nil {
			return nil, err
		}
		return res.Count, nil
	})
	if err != nil {
		return nil, err
	}

	return &api.StreamCountResponse{
		Count: item.Value().(int64),
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

	prefix := BUCKET_STREAM_LIST

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

		var stream uuid.UUID
		if err := stream.UnmarshalBinary(cursor.Item().Key()[1:]); err != nil {
			return nil, err
		}

		res.Streams = append(res.Streams, stream.String())
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

	cache := ccache.New(ccache.Configure())

	store := &BadgerEventStore{db, cache}

	if !db.Opts().InMemory {
		go func() {
			if err := db.RunValueLogGC(0.7); err != nil && err != badger.ErrNoRewrite {
				log.Fatal(err)
			}

			time.Sleep(time.Second)
		}()
	}

	return store, nil
}

func getStreamVersionKey(stream uuid.UUID, version uint32) []byte {
	data := make([]byte, 8)
	binary.BigEndian.PutUint32(data, version)

	result := BUCKET_STREAMS
	result = append(result, stream[:]...)
	result = append(result, data...)
	return result
}

func getStreamKey(stream uuid.UUID) []byte {
	result := BUCKET_STREAMS
	result = append(result, stream[:]...)
	return result
}

func getStreamListKey(stream uuid.UUID) []byte {
	result := BUCKET_STREAM_LIST
	result = append(result, stream[:]...)
	return result
}

func getEventKey(id ulid.ULID) []byte {
	result := BUCKET_EVENTS
	result = append(result, id[:]...)
	return result
}

func getCurrentSequence(txn *badger.Txn) (sequence uint64, err error) {
	item, err := txn.Get(KEY_CURRENT_SEQUENCE)
	if err == badger.ErrKeyNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	err = item.Value(func(val []byte) error {
		sequence = binary.BigEndian.Uint64(val)

		return nil
	})

	return sequence, err
}

func setCurrentSequence(txn *badger.Txn, sequence uint64) error {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, sequence)
	return txn.Set(KEY_CURRENT_SEQUENCE, data)
}

func getSequenceKey(sequence uint64) []byte {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, sequence)

	result := BUCKET_SEQUENCE
	result = append(result, data...)
	return result
}
