package fsm

import (
	"io"

	"github.com/dgraph-io/badger"
	"github.com/hashicorp/raft"
)

type badgerFSM struct {
	db *badger.DB
}

func (b *badgerFSM) Apply(log *raft.Log) interface{} {
	return nil
}

func (b *badgerFSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (b *badgerFSM) Restore(io.ReadCloser) error {
	return nil
}

func NewBadgerFSM(db *badger.DB) *badgerFSM {
	return &badgerFSM{db}
}
