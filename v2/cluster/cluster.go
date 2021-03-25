package cluster

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

type Cluster struct {
	raft *raft.Raft
}

func NewCluster(localID string, bindAddr string, advrAddr string, fsm raft.FSM) (*Cluster, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID)

	lst, err := raftboltdb.NewBoltStore("data/log.dat")
	if err != nil {
		return nil, fmt.Errorf("Failed to create Raft log store: %v", err)
	}

	sst, err := raftboltdb.NewBoltStore("data/stable.dat")
	if err != nil {
		return nil, fmt.Errorf("Failed to create Raft stable store: %v", err)
	}

	fss, err := raft.NewFileSnapshotStore("data", 1, os.Stdout)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Raft snapshot store: %v", err)
	}

	addr, err := net.ResolveTCPAddr("tcp", advrAddr)
	if err != nil {
		return nil, fmt.Errorf("Failed to resolve address: %v", err)
	}

	transport, err := raft.NewTCPTransport(bindAddr, addr, 10, time.Second*5, os.Stdout)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Raft transport: %v", err)
	}

	raftServer, err := raft.NewRaft(config, fsm, lst, sst, fss, transport)
	if err != nil {
		if err != nil {
			return nil, fmt.Errorf("Failed to create Raft: %v", err)
		}
	}

	return &Cluster{raftServer}, nil
}
