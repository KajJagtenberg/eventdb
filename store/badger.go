package store

import (
	"bytes"
	"errors"

	"github.com/dgraph-io/badger/v3"
)

type BadgerEventStore struct {
	db *badger.DB
}

var (
	MAGIC_NUMBER = []byte{32, 179}

	BUCKET_EVENTS   = []byte{0, 0}
	BUCKET_STREAMS  = []byte{0, 1}
	BUCKET_METADATA = []byte{0, 2}
)

func NewBadgerEventStore(db *badger.DB) (*BadgerEventStore, error) {
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

	return &BadgerEventStore{db}, nil
}
