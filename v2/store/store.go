package store

import (
	"math/rand"

	"github.com/oklog/ulid"
	"go.etcd.io/bbolt"
)

var (
	entropy = ulid.Monotonic(rand.New(rand.NewSource(int64(ulid.Now()))), 0)
)

type Storage struct {
	db *bbolt.DB
}

func (s *Storage) Add(req *AddRequest) error {
	return nil
}

func (s *Storage) Get(req *GetRequest) error {
	return nil
}

func (s *Storage) Log(req *LogRequest) error {
	return nil
}

func NewStorage(db *bbolt.DB) *Storage {
	return &Storage{db}
}
