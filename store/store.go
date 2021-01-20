package store

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"github.com/vmihailenco/msgpack"
)

var (
	PrefixEvent  = []byte{0, 1}
	PrefixStream = []byte{0, 2}

	entropy = ulid.Monotonic(rand.New((rand.NewSource((int64(ulid.Now()))))), 0)
)

type Store struct {
	db *badger.DB
}

func (s *Store) AppendToStream(streamId uuid.UUID, version int, events []AppendEvent) error {
	return s.db.Update(func(txn *badger.Txn) error {
		stream, err := getStream(txn, streamId)
		if err != nil {
			return err
		}

		if version != stream.Version {
			return fmt.Errorf("Optimistic concurrency error: given = %d, current = %d", version, stream.Version)
		}

		for i, insert := range events {
			id, err := ulid.New(ulid.Now(), entropy)
			if err != nil {
				return err
			}

			event := Event{
				ID:            id,
				Stream:        streamId,
				Version:       version + i,
				Type:          insert.Type,
				Data:          insert.Data,
				Metadata:      insert.Metadata,
				CausationID:   insert.CausationID,
				CorrelationID: insert.CorrelationID,
				Timestamp:     time.Now(),
			}

			marshalledID, err := id.MarshalBinary()
			if err != nil {
				return err
			}

			eventKey := append(PrefixEvent, marshalledID...)

			value, err := msgpack.Marshal(event)
			if err != nil {
				return err
			}

			if err := txn.Set(eventKey, value); err != nil {
				return err
			}

			stream.Events = append(stream.Events, id)
			stream.Version++
		}

		return setStream(txn, streamId, stream)
	})
}

func (s *Store) LoadFromStream(streamId uuid.UUID, version int, limit int) (int, []Event, error) {
	result := []Event{}
	streamVersion := 0

	err := s.db.View(func(txn *badger.Txn) error {
		stream, err := getStream(txn, streamId)

		if err != nil {
			return err
		}

		streamVersion = stream.Version

		for _, ref := range stream.Events {
			if version > 0 {
				version--
				continue
			}

			item, err := txn.Get(append(PrefixEvent, ref[:]...))
			if err != nil {
				return err
			}

			value, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			var event Event

			if err := msgpack.Unmarshal(value, &event); err != nil {
				return err
			}

			result = append(result, event)
		}

		return nil
	})

	if err != nil {
		return 0, nil, err
	}

	return streamVersion, result, nil
}

func (s *Store) GetStreams(offset int, limit int) ([]uuid.UUID, error) {
	result := []uuid.UUID{}

	if limit > 100 {
		limit = 100
	}

	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := PrefixStream
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := item.Key()

			if len(key) > 18 {
				continue
			}

			if offset > 0 {
				offset--
				continue
			}

			if len(result) >= limit {
				return nil
			}

			stream, err := uuid.FromBytes(key[2:])
			if err != nil {
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

func (s *Store) Backup(dst io.Writer) error {
	if _, err := s.db.Backup(dst, 0); err != nil {
		return err
	}

	return nil
}

func NewStore(db *badger.DB) *Store {
	return &Store{db}
}

func getStream(txn *badger.Txn, stream uuid.UUID) (*Stream, error) {
	key := append(PrefixStream, stream[:]...)

	item, err := txn.Get(key)
	if err == badger.ErrKeyNotFound {
		return &Stream{}, nil
	} else if err != nil {
		return nil, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	var result *Stream

	if err := msgpack.Unmarshal(value, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func setStream(txn *badger.Txn, stream uuid.UUID, value *Stream) error {
	key := append(PrefixStream, stream[:]...)

	data, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}

	return txn.Set(key, data)
}

// func getStreamKey(stream uuid.UUID, version int) []byte {
// 	indexVersion := make([]byte, 4)
// 	binary.BigEndian.PutUint32(indexVersion, uint32(version))
// 	streamKey := append([]byte(stream[:]), indexVersion...)
// 	streamKey = append(PrefixStream, streamKey...)
// 	return streamKey
// }

// func getStreamVersion(txn *badger.Txn, stream uuid.UUID) (int, error) {
// 	key := append(PrefixStreamVersion, stream[:]...)

// 	item, err := txn.Get(key)
// 	if err == badger.ErrKeyNotFound {
// 		return 0, nil
// 	} else if err != nil {
// 		return 0, err
// 	}

// 	value, err := item.ValueCopy(nil)
// 	if err != nil {
// 		return 0, err
// 	}

// 	version := binary.BigEndian.Uint32(value)

// 	return int(version + 1), nil
// }

// func setStreamVersion(txn *badger.Txn, stream uuid.UUID, version int) error {
// 	value := make([]byte, 4)
// 	binary.LittleEndian.PutUint32(value, uint32(version))

// 	key := append(PrefixStreamVersion, stream[:]...)

// 	return txn.Set(key, value)
// }
