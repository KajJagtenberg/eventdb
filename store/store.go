package store

import (
	"io"
	"math/rand"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"go.etcd.io/bbolt"
)

var (
	PrefixEvent  = []byte{0, 1}
	PrefixStream = []byte{0, 2}

	entropy = ulid.Monotonic(rand.New((rand.NewSource((int64(ulid.Now()))))), 0)
)

type Store struct {
	db *bbolt.DB
}

func (s *Store) AppendToStream(streamId uuid.UUID, version int, events []AppendEvent) error {
	return nil
}

func (s *Store) LoadFromStream(streamId uuid.UUID, version int, limit int) ([]Event, error) {
	return nil, nil
}

func (s *Store) GetStreams(offset int, limit int) ([]uuid.UUID, error) {
	return nil, nil
}

func (s *Store) Backup(dst io.Writer) error {
	return nil
}

func NewStore(db *bbolt.DB) *Store {
	return &Store{db}
}
