package store

import "go.etcd.io/bbolt"

var (
	buckets = []string{"streams", "events"}
)

type Store struct {
	db *bbolt.DB
}

func (s *Store) Size() int64 {
	var size int64 = 0

	s.db.View(func(t *bbolt.Tx) error {
		size = t.Size()

		return nil
	})

	return size
}

func NewStore(db *bbolt.DB) (*Store, error) {
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

	return &Store{db}, nil
}
