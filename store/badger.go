package store

import (
	"encoding/binary"
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
	db          *badger.DB
	systemCache *ccache.Cache
	eventCache  *ccache.Cache
}

var (
	BUCKET_EVENTS      = []byte{0}
	BUCKET_STREAMS     = []byte{1}
	BUCKET_SEQUENCE    = []byte{2}
	BUCKET_SYSTEM      = []byte{3}
	BUCKET_STREAM_LIST = []byte{4}

	KEY_CURRENT_SEQUENCE = []byte{3, 0}
)

func (s *BadgerEventStore) AppendToStream(req *api.AppendToStreamRequest) (res *api.AppendToStreamResponse, err error) {
	res = &api.AppendToStreamResponse{
		Events: make([]string, 0),
	}

	// stream, err := uuid.Parse(req.Stream)
	// if err != nil {
	// 	return nil, err
	// }

	if len(req.Events) == 0 {
		return nil, ErrEmptyEvents
	}

	if req.Version < -1 {
		return nil, ErrWrongVersion
	}

	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	return res, txn.Commit()
}

func (s *BadgerEventStore) GetStream(req *api.GetStreamRequest) (res *api.GetStreamResponse, err error) {
	return nil, nil
}

func (s *BadgerEventStore) GetGlobalStream(req *api.GetGlobalStreamRequest) (res *api.GetGlobalStreamResponse, err error) {
	return nil, nil
}

func (s *BadgerEventStore) GetEvent(req *api.GetEventRequest) (res *api.Event, err error) {
	return nil, nil
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
	item, err := s.systemCache.Fetch("EVENT_COUNT", ESTIMATE_TTL, func() (interface{}, error) {
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
	item, err := s.systemCache.Fetch("STREAM_COUNT", ESTIMATE_TTL, func() (interface{}, error) {
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

	systemCache := ccache.New(ccache.Configure())
	eventCache := ccache.New(ccache.Configure())

	store := &BadgerEventStore{
		db:          db,
		systemCache: systemCache,
		eventCache:  eventCache,
	}

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

func getEvent(txn *badger.Txn, cache *ccache.Cache, id ulid.ULID) (*api.Event, error) {
	item, err := cache.Fetch(id.String(), time.Hour, func() (interface{}, error) {
		item, err := txn.Get(getEventKey(id))
		if err != nil {
			return nil, err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}

		var event PersistedEvent
		if err := proto.Unmarshal(val, &event); err != nil {
			return nil, err
		}

		var id ulid.ULID
		if err := id.UnmarshalBinary(event.Id); err != nil {
			return nil, err
		}

		var stream uuid.UUID
		if err := stream.UnmarshalBinary(event.Stream); err != nil {
			return nil, err
		}

		var causationID ulid.ULID
		if err := causationID.UnmarshalBinary(event.CausationId); err != nil {
			return nil, err
		}

		var correlationID ulid.ULID
		if err := correlationID.UnmarshalBinary(event.CorrelationId); err != nil {
			return nil, err
		}

		return &api.Event{
			Id:            id.String(),
			Stream:        stream.String(),
			Version:       event.Version,
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   causationID.String(),
			CorrelationId: correlationID.String(),
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return item.Value().(*api.Event), nil
}
