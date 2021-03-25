package cluster

import (
	"net"
	"os"
	"time"

	"github.com/hashicorp/raft"
)

type Cluster struct {
	raft *raft.Raft
}

func NewCluster(localID string, bindAddr string, advrAddr string, fsm raft.FSM, lst raft.LogStore, sst raft.StableStore, sns raft.SnapshotStore) (*Cluster, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID)

	addr, err := net.ResolveTCPAddr("tcp", advrAddr)
	if err != nil {
		return nil, err
	}

	transport, err := raft.NewTCPTransport(bindAddr, addr, 10, time.Second*5, os.Stdout)

	raftServer, err := raft.NewRaft(config, fsm, lst, sst, sns, transport)
	if err != nil {
		return nil, err
	}

	return &Cluster{raftServer}, nil
}
