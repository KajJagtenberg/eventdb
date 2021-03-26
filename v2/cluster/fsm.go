package cluster

import (
	"errors"
	"io"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/raft"
	"go.etcd.io/bbolt"
)

type FSM struct {
	db *bbolt.DB
}

func (fsm *FSM) Apply(applyLog *raft.Log) interface{} {
	switch applyLog.Type {
	case raft.LogCommand:
		var cmd ApplyLog

		if err := proto.Unmarshal(applyLog.Data, &cmd); err != nil {
			log.Printf("Failed to unmarshal ApplyLog: %v", err)
			return err
		}

		switch cmd := cmd.Command.(type) {
		case *ApplyLog_Add:
			log.Printf("[FSM] Add command received: %v", cmd.Add)
		}
	default:
		log.Println("Type:", applyLog.Type)
	}

	return nil
}

func (fsm *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, errors.New("Not implemented")
}

func (fsm *FSM) Restore(io.ReadCloser) error {
	return errors.New("Not implemented")
}

func NewFSM(db *bbolt.DB) (*FSM, error) {
	err := db.Update(func(t *bbolt.Tx) error {
		buckets := []string{"streams", "events"}

		for _, bucket := range buckets {
			if _, err := t.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &FSM{db}, nil
}
