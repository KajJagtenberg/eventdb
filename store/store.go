package store

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/oklog/ulid"
)

var (
	PrefixEvent  = []byte{0, 1}
	PrefixStream = []byte{0, 2}

	entropy = ulid.Monotonic(rand.New((rand.NewSource((int64(ulid.Now()))))), 0)
)

type Store struct {
	db *badger.DB
}

func (s *Store) AppendToStream(stream string, version int, events []AppendEvent) error {
	return s.db.Update(func(txn *badger.Txn) error {
		streamKey := getStreamKey(stream, version)

		_, err := txn.Get(streamKey)
		if err == nil {
			return errors.New("Concurrent stream modification")
		} else if err != badger.ErrKeyNotFound {
			return err
		}

		for i, insert := range events {
			// TODO: Check if the version already exists and return an error if it does

			id, err := ulid.New(ulid.Now(), entropy)
			if err != nil {
				return err
			}

			event := Event{
				ID:            id.String(),
				Stream:        stream,
				Version:       version + i,
				Type:          insert.Type,
				Data:          insert.Data,
				CausationID:   insert.CausationID,
				CorrelationID: insert.CorrelationID,
				Timestamp:     time.Now(),
			}

			marshalledID, err := id.MarshalBinary()
			if err != nil {
				return err
			}

			eventKey := append(PrefixEvent, marshalledID...)

			value, err := json.Marshal(event)
			if err != nil {
				return err
			}

			if err := txn.Set(eventKey, value); err != nil {
				return err
			}

			streamKey := getStreamKey(stream, version+i)

			if err := txn.Set(streamKey, marshalledID); err != nil {
				return err
			}
		}

		return nil
	})
}

func getStreamKey(stream string, version int) []byte {
	indexVersion := make([]byte, 4)
	binary.BigEndian.PutUint32(indexVersion, uint32(version))
	streamKey := append([]byte(stream), indexVersion...)
	streamKey = append(PrefixStream, streamKey...)
	return streamKey
}

func (s *Store) LoadFromStream(stream string, version int, limit int) ([]Event, error) {
	result := []Event{}

	err := s.db.View(func(txn *badger.Txn) error {
		for i := 0; len(result) < limit || limit == 0; i++ {
			streamKey := getStreamKey(stream, version+i)

			item, err := txn.Get(streamKey)
			if err == badger.ErrKeyNotFound {
				return nil
			} else if err != nil {
				return err
			}

			id, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			item, err = txn.Get(append(PrefixEvent, id...))
			if err != nil {
				return err
			}

			value, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			var event Event
			if err := json.Unmarshal(value, &event); err != nil {
				return err
			}

			result = append(result, event)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewStore(db *badger.DB) *Store {
	return &Store{db}
}
