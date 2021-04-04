package store

import (
	"io"

	"go.etcd.io/bbolt"
)

var (
	buckets = []string{"streams", "events"}
)

type BoltStore struct {
	db *bbolt.DB
}

func (s *BoltStore) Size() int64 {
	var size int64 = 0

	s.db.View(func(t *bbolt.Tx) error {
		size = t.Size()

		return nil
	})

	return size
}

func (s *BoltStore) Backup(dst io.Writer) error {
	return s.db.View(func(t *bbolt.Tx) error {
		if _, err := t.WriteTo(dst); err != nil {
			return err
		}

		return nil
	})
}

func NewBoltStore(db *bbolt.DB) (*BoltStore, error) {
	if err := db.Update(func(t *bbolt.Tx) error {
		for _, bucket := range buckets {
			if _, err := t.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &BoltStore{db}, nil
}
