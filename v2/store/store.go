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

func NewStorage(db *bbolt.DB) *Storage {
	return &Storage{db}
}
