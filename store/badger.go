package store

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"hash/crc32"
	"io"
	"log"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

type BadgerEventStore struct {
	db                  *badger.DB
	estimateStreamCount int64
	estimateEventCount  int64
}

var (
	MAGIC_NUMBER = []byte{32, 179}

	BUCKET_EVENTS   = []byte{0, 0}
	BUCKET_STREAMS  = []byte{0, 1}
	BUCKET_METADATA = []byte{0, 2}
)

var (
	ErrConcurrentStreamModification = errors.New("concurrent stream modification")
	ErrGappedStream                 = errors.New("given version leaves gap in stream")
)

func (s *BadgerEventStore) Size() (int64, error) {
	lsm, vlog := s.db.Size()

	return lsm + vlog, nil
}

func (s *BadgerEventStore) Backup(dst io.Writer) error {
	_, err := s.db.Backup(dst, 0)
	return err
}

func (s *BadgerEventStore) Add(stream uuid.UUID, version uint32, events []EventData) ([]Event, error) {
	if bytes.Equal(stream[:], make([]byte, 16)) {
		return nil, errors.New("stream cannot be all zeroes")
	}

	if len(events) == 0 {
		return nil, errors.New("list of events is empty")
	}

	result := make([]Event, 0)

	if err := s.db.Update(func(txn *badger.Txn) error {
		var s Stream

		item, err := txn.Get(getStreamKey(stream))
		if err == nil {
			if err := item.Value(func(val []byte) error {
				return s.Unmarshal(val)
			}); err != nil {
				return err
			}
		} else if err == badger.ErrKeyNotFound {
			s.ID = stream
			s.AddedAt = time.Now()
		} else {
			return err
		}

		if int(version) < len(s.Events) {
			return ErrConcurrentStreamModification
		}

		if int(version) > len(s.Events) {
			return ErrGappedStream
		}

		now := time.Now()

		for i, event := range events {
			if len(event.Type) == 0 {
				return errors.New("event type cannot be empty")
			}

			// if len(event.Data) == 0 {
			// 	return errors.New("event data cannot be empty")
			// }

			id, err := ulid.New(ulid.Now(), rand.Reader)
			if err != nil {
				return err
			}

			record := Event{
				ID:            id,
				Stream:        stream,
				Version:       version + uint32(i),
				Type:          event.Type,
				Data:          event.Data,
				Metadata:      event.Metadata,
				CausationID:   event.CausationID,
				CorrelationID: event.CorrelationID,
				AddedAt:       now,
			}

			if bytes.Equal(record.CausationID[:], make([]byte, 16)) {
				record.CausationID = record.ID
			}

			if bytes.Equal(record.CorrelationID[:], make([]byte, 16)) {
				record.CorrelationID = record.CausationID
			}

			result = append(result, record)

			s.Events = append(s.Events, record.ID)
		}

		for _, record := range result {
			if value, err := record.Marshal(); err != nil {
				return err
			} else {
				if err := txn.Set(getEventKey(record.ID), value); err != nil {
					return err
				}
			}
		}

		if value, err := s.Marshal(); err != nil {
			return err
		} else {
			if err := txn.Set(getStreamKey(stream), value); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *BadgerEventStore) Get(stream uuid.UUID, version uint32, limit uint32) ([]Event, error) {
	result := make([]Event, 0)

	if err := s.db.View(func(txn *badger.Txn) error {
		var s Stream

		item, err := txn.Get(getStreamKey(stream))
		if err == nil {
			if err := item.Value(func(val []byte) error {
				return s.Unmarshal(val)
			}); err != nil {
				return err
			}
		} else if err == badger.ErrKeyNotFound {
			return nil
		} else {
			return err
		}

		for _, id := range s.Events {
			if version > 0 {
				version--
				continue
			}
			if len(result) >= int(limit) && limit != 0 {
				return nil
			}
			item, err := txn.Get(getEventKey(id))
			if err == nil {

			} else if err == badger.ErrKeyNotFound {
				return errors.New("event cannot be found. this should never happen")
			} else {
				return err
			}

			var event Event

			if err := item.Value(func(val []byte) error {
				return event.Unmarshal(val)
			}); err != nil {
				return err
			}

			result = append(result, event)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *BadgerEventStore) GetAll(offset ulid.ULID, limit uint32) ([]Event, error) {
	if limit == 0 {
		limit = 10
	}

	result := make([]Event, 0)

	if err := s.db.View(func(txn *badger.Txn) error {
		cursor := txn.NewIterator(badger.DefaultIteratorOptions)
		defer cursor.Close()

		for cursor.Seek(getEventKey(offset)); cursor.ValidForPrefix(BUCKET_EVENTS); cursor.Next() {
			if bytes.Equal(cursor.Item().Key(), getEventKey(offset)) {
				continue
			}

			if len(result) >= int(limit) {
				break
			}

			var event Event

			if err := cursor.Item().Value(func(val []byte) error {
				return event.Unmarshal(val)
			}); err != nil {
				return err
			}

			result = append(result, event)

			// cursor.Next()
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *BadgerEventStore) EventCount() (int64, error) {
	var total int64

	if err := s.db.View(func(txn *badger.Txn) error {
		cursor := txn.NewIterator(badger.DefaultIteratorOptions)

		defer cursor.Close()

		for cursor.Seek(BUCKET_EVENTS); cursor.ValidForPrefix(BUCKET_EVENTS); cursor.Next() {
			total++
		}

		return nil
	}); err != nil {
		return 0, err
	}

	return total, nil
}

func (s *BadgerEventStore) StreamCount() (int64, error) {
	var total int64

	if err := s.db.View(func(txn *badger.Txn) error {
		cursor := txn.NewIterator(badger.DefaultIteratorOptions)

		defer cursor.Close()

		for cursor.Seek(BUCKET_STREAMS); cursor.ValidForPrefix(BUCKET_STREAMS); cursor.Next() {
			total++
		}

		return nil
	}); err != nil {
		return 0, err
	}

	return total, nil
}

func (s *BadgerEventStore) EventCountEstimate() (int64, error) {
	return s.estimateEventCount, nil
}

func (s *BadgerEventStore) StreamCountEstimate() (int64, error) {
	return s.estimateStreamCount, nil
}

func (s *BadgerEventStore) ListStreams(skip uint32, limit uint32) ([]Stream, error) {
	result := make([]Stream, 0)

	if limit == 0 {
		limit = 25
	}

	err := s.db.View(func(txn *badger.Txn) error {
		prefix := BUCKET_STREAMS

		cursor := txn.NewIterator(badger.DefaultIteratorOptions)

		defer cursor.Close()

		for cursor.Seek(prefix); cursor.ValidForPrefix(prefix); cursor.Next() {
			if skip > 0 {
				skip--
				continue
			}

			if len(result) >= int(limit) {
				return nil
			}

			var stream Stream

			if err := cursor.Item().Value(func(val []byte) error {
				return stream.Unmarshal(val)
			}); err != nil {
				return err
			}

			result = append(result, stream)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *BadgerEventStore) Checksum() (ulid.ULID, []byte, error) {
	h := crc32.NewIEEE()
	checksum := Checksum{
		ID:  ulid.ULID{},
		Sum: make([]byte, 4),
	}

	err := s.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(getMetadataKey([]byte("checksum")))
		if err == nil {
			if err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &checksum)
			}); err != nil {
				return err
			}
		} else if err != badger.ErrKeyNotFound {
			return err
		}

		cursor := txn.NewIterator(badger.DefaultIteratorOptions)
		defer cursor.Close()

		for cursor.Seek(getEventKey(checksum.ID)); cursor.ValidForPrefix(BUCKET_EVENTS); cursor.Next() {
			item := cursor.Item()

			if bytes.Equal(item.Key(), getEventKey(checksum.ID)) {
				return nil
			}

			h.Reset()

			if err := checksum.ID.UnmarshalBinary(item.Key()[2:]); err != nil {
				return err
			}

			if err := item.Value(func(val []byte) error {
				_, err := h.Write(val)
				return err
			}); err != nil {
				return err
			}

			if _, err := h.Write(checksum.Sum); err != nil {
				return err
			}

			checksum.Sum = h.Sum(nil)
		}

		if value, err := json.Marshal(checksum); err != nil {
			return err
		} else {
			if err := txn.Set(getMetadataKey([]byte("checksum")), value); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return checksum.ID, checksum.Sum, err
	}

	return checksum.ID, checksum.Sum, nil
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

	if err := db.Update(func(txn *badger.Txn) error {
		k := append(BUCKET_METADATA, []byte("MAGIC_NUMBER")...)

		item, err := txn.Get(k)
		if err == nil {
			if err := item.Value(func(val []byte) error {
				if !bytes.Equal(val, MAGIC_NUMBER) {
					return errors.New("invalid magic number found. database could be in a corrupt state")
				}

				return nil
			}); err != nil {
				return err
			}
		} else if err == badger.ErrKeyNotFound {
			if err := txn.Set(k, MAGIC_NUMBER); err != nil {
				return err
			}
		} else {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	store := &BadgerEventStore{db, 0, 0}

	if options.EstimateCounts {
		go func() {
			for {
				streamCount, err := store.StreamCount()
				if err != nil {
					log.Fatalf("failed to get stream count: %v", err)
				}

				eventCount, err := store.EventCount()
				if err != nil {
					log.Fatalf("failed to get stream count: %v", err)
				}

				store.estimateStreamCount = streamCount
				store.estimateEventCount = eventCount

				time.Sleep(ESTIMATE_SLEEP_TIME)
			}
		}()
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
		again:
			err := db.RunValueLogGC(0.7)
			if err == nil {
				goto again
			}
		}
	}()

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

func getMetadataKey(key []byte) []byte {
	result := BUCKET_METADATA
	result = append(result, key...)
	return result
}

type Checksum struct {
	ID  ulid.ULID `json:"ulid"`
	Sum []byte    `json:"sum"`
}
