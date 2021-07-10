package store

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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
	PREFIX_EVENT  = []byte{0}
	PREFIX_STREAM = []byte{1}
	SEPERATOR     = []byte{58}
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

	events, err := convertToEvents(req)
	if err != nil {
		return nil, err
	}

	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	for i, event := range events {
		data, err := json.Marshal(event)
		if err != nil {
			return nil, err
		}

		value, err := json.Marshal(&Value{
			Encoding: 0,
			Data:     data,
			Version:  0,
		})
		if err != nil {
			return nil, err
		}

		if err := txn.Set(getEventKey(event.Id), value); err != nil {
			return nil, err
		}

		if err := txn.Set(getStreamKey(stream, int(req.Version)+i), nil); err != nil {
			return nil, err
		}

		res.Events = append(res.Events, event.Id)
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

func getStreamKey(stream uuid.UUID, version int) []byte {
	buf := new(bytes.Buffer)
	buf.Write(PREFIX_STREAM)
	buf.Write(SEPERATOR)
	buf.WriteString(stream.String())
	buf.Write(SEPERATOR)

	binary.Write(buf, binary.BigEndian, uint32(version))

	return buf.Bytes()
}

func getEventKey(id string) []byte {
	buf := new(bytes.Buffer)
	buf.Write(PREFIX_EVENT)
	buf.Write(SEPERATOR)
	buf.WriteString(id)
	buf.Write(SEPERATOR)

	return buf.Bytes()
}
