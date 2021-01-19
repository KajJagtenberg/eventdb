package store

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
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

func (s *Store) AppendToStream(stream uuid.UUID, version int, events []AppendEvent) error {
	return s.db.Update(func(txn *badger.Txn) error {
		streamVersion, err := getStreamVersion(txn, stream)
		if err != nil {
			return err
		}

		if version < streamVersion {
			return errors.New("Concurrent stream modification")
		}
		if version > streamVersion {
			return errors.New("Version does not line up with stream version")
		}

		for i, insert := range events {
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

			if err := setStreamVersion(txn, stream, version+i); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Store) LoadFromStream(stream uuid.UUID, version int, limit int) ([]Event, error) {
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

func getStreamKey(stream uuid.UUID, version int) []byte {
	indexVersion := make([]byte, 4)
	binary.BigEndian.PutUint32(indexVersion, uint32(version))
	streamKey := append([]byte(stream[:]), indexVersion...)
	streamKey = append(PrefixStream, streamKey...)
	return streamKey
}

func getStreamVersion(txn *badger.Txn, stream uuid.UUID) (int, error) {
	key := append(PrefixStream, stream[:]...)

	item, err := txn.Get(key)
	if err == badger.ErrKeyNotFound {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return 0, err
	}

	version := binary.BigEndian.Uint32(value)

	return int(version + 1), nil
}

func setStreamVersion(txn *badger.Txn, stream uuid.UUID, version int) error {
	value := make([]byte, 4)
	binary.LittleEndian.PutUint32(value, uint32(version))

	key := append(PrefixStream, stream[:]...)

	return txn.Set(key, value)
}
