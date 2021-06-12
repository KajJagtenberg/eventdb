package fsm

import (
	"errors"
	"io"

	"github.com/hashicorp/raft"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/store"
	"google.golang.org/protobuf/proto"
)

type fsm struct {
	store store.EventStore
}

func (b *fsm) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:
		var cmd api.Command
		if err := proto.Unmarshal(log.Data, &cmd); err != nil {
			return &ApplyResponse{
				Error: err,
			}
		}

		switch cmd.Op {
		case "ADD":
			var req *api.AddRequest
			if err := proto.Unmarshal(cmd.Payload, req); err != nil {
				return &ApplyResponse{
					Error: err,
				}
			}
			res, err := b.store.Add(req)
			return &ApplyResponse{
				Data:  res,
				Error: err,
			}
		}
	}

	return &ApplyResponse{
		Error: errors.New("unknown command"),
	}
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
