package fsm

import (
	"io"

	"github.com/hashicorp/raft"
	"github.com/kajjagtenberg/eventflowdb/store"
)

type fsm struct {
	store store.EventStore
}

func (b *fsm) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:

	}

	return nil
}

func (b *fsm) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (b *fsm) Restore(io.ReadCloser) error {
	return nil
}
func NewFSM(store store.EventStore) *fsm {
	return &fsm{store}
}
