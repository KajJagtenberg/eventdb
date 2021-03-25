package cluster

import "github.com/hashicorp/raft"

type FSM struct{}

func NewFSM() (raft.FSM, error) {
	return &raft.MockFSM{}, nil
}
