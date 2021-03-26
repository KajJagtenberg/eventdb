package cluster

import (
	"io"
	"log"

	"github.com/hashicorp/raft"
)

type FSM struct{}

func (fsm *FSM) Apply(applyLog *raft.Log) interface{} {
	switch applyLog.Type {
	case raft.LogCommand:
		log.Println("Log command")
	default:
		log.Println("Type:", applyLog.Type)
	}

	return nil
}

func (fsm *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (fsm *FSM) Restore(io.ReadCloser) error {
	return nil
}

func NewFSM() (raft.FSM, error) {
	return &FSM{}, nil
}
