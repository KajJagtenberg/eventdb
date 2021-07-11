package store

import (
	"io"
	"log"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/eventflowdb/eventflowdb/api"
	"github.com/eventflowdb/eventflowdb/conv"
	"github.com/google/uuid"
)

type BadgerEventStore struct {
	db *badger.DB
}

var (
	PREFIX_EVENT         = []byte{0}
	PREFIX_STREAM        = []byte{1}
	PREFIX_STREAM_EVENT  = []byte{2}
	PREFIX_GLOBAL_STREAM = []byte{3}
	SEPERATOR            = []byte{58}
)

func (s *BadgerEventStore) AppendToStream(req *api.AppendToStreamRequest) (res *api.AppendToStreamResponse, err error) {
	res = &api.AppendToStreamResponse{
		Events: make([]string, 0),
	}

	stream, err := uuid.Parse(req.Stream)
	if err != nil {
		return nil, err
	}

	if len(req.Events) == 0 {
		return nil, ErrEmptyEvents
	}

	if req.Version < -1 {
		return nil, ErrWrongVersion
	}

	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	persistedStream, err := getStream(txn, stream)
	if err != nil {
		return nil, err
	}

	if req.Version != -1 && persistedStream.Version != uint32(req.Version) {
		return nil, ErrConcurrentStreamModification
	}

	for i, event := range req.Events {
		if len(event.Type) == 0 {
			return nil, ErrEmptyEventType
		}

		id := uuid.New()

		if len(event.CausationId) == 0 {
			event.CausationId = id.String()
		}
		if len(event.CorrelationId) == 0 {
			event.CorrelationId = id.String()
		}

		version := persistedStream.Version + uint32(i)

		if err := setEvent(txn, id, stream, event, version); err != nil {
			return nil, err
		}

		if err := setStreamEvent(txn, stream, version, id); err != nil {
			return nil, err
		}

		res.Events = append(res.Events, id.String())
	}

	persistedStream.Version += uint32(len(req.Events))

	if err := setStream(txn, stream, persistedStream); err != nil {
		return nil, err
	}

	if err := txn.Commit(); err != nil {
		return nil, err
	}

	return res, nil
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
	return nil, nil
}

func (s *BadgerEventStore) StreamCount(req *api.StreamCountRequest) (res *api.StreamCountResponse, err error) {
	return nil, nil
}

func (s *BadgerEventStore) Size(req *api.SizeRequest) (res *api.SizeResponse, err error) {
	res = &api.SizeResponse{}

	lsm, _ := s.db.Size()

	res.Size = lsm
	res.SizeHuman = conv.ByteCountSI(res.Size)

	return res, nil
}

func (s *BadgerEventStore) ListStreams(req *api.ListStreamsRequest) (res *api.ListStreamsReponse, err error) {
	return nil, nil
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

	store := &BadgerEventStore{
		db: db,
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
