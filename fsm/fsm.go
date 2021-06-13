package fsm

import (
	"errors"
	"io"
	"log"

	"github.com/hashicorp/raft"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/store"
	"google.golang.org/protobuf/proto"
)

type fsm struct {
	store store.EventStore
}

func (b *fsm) Apply(l *raft.Log) interface{} {
	switch l.Type {
	case raft.LogCommand:
		var cmd api.Command
		if err := proto.Unmarshal(l.Data, &cmd); err != nil {
			return &ApplyResponse{
				Error: err,
			}
		}

		// log.Println(cmd.Op)

		switch cmd.Op {
		case "ADD":
			var req api.AddRequest
			if err := proto.Unmarshal(cmd.Payload, &req); err != nil {
				return &ApplyResponse{
					Error: err,
				}
			}
			res, err := b.store.Add(&req)
			return &ApplyResponse{
				Data:  res,
				Error: err,
			}

		case "GET":
			var req api.GetRequest
			if err := proto.Unmarshal(cmd.Payload, &req); err != nil {
				return &ApplyResponse{
					Error: err,
				}
			}
			res, err := b.store.Get(&req)
			return &ApplyResponse{
				Data:  res,
				Error: err,
			}

		case "GET_ALL":
			var req api.GetAllRequest
			if err := proto.Unmarshal(cmd.Payload, &req); err != nil {
				return &ApplyResponse{
					Error: err,
				}
			}
			res, err := b.store.GetAll(&req)
			return &ApplyResponse{
				Data:  res,
				Error: err,
			}

		case "EVENT_COUNT":
			var req api.EventCountRequest
			if err := proto.Unmarshal(cmd.Payload, &req); err != nil {
				return &ApplyResponse{
					Error: err,
				}
			}
			res, err := b.store.EventCount(&req)
			return &ApplyResponse{
				Data:  res,
				Error: err,
			}

		case "EVENT_COUNT_ESTIMATE":
			var req api.EventCountEstimateRequest
			if err := proto.Unmarshal(cmd.Payload, &req); err != nil {
				return &ApplyResponse{
					Error: err,
				}
			}
			res, err := b.store.EventCountEstimate(&req)
			return &ApplyResponse{
				Data:  res,
				Error: err,
			}

		case "STREAM_COUNT":
			var req api.StreamCountRequest
			if err := proto.Unmarshal(cmd.Payload, &req); err != nil {
				return &ApplyResponse{
					Error: err,
				}
			}
			res, err := b.store.StreamCount(&req)
			return &ApplyResponse{
				Data:  res,
				Error: err,
			}

		case "STREAM_COUNT_ESTIMATE":
			var req api.StreamCountEstimateRequest
			if err := proto.Unmarshal(cmd.Payload, &req); err != nil {
				return &ApplyResponse{
					Error: err,
				}
			}
			res, err := b.store.StreamCountEstimate(&req)
			return &ApplyResponse{
				Data:  res,
				Error: err,
			}

		case "LIST_STREAMS":
			var req api.ListStreamsRequest
			if err := proto.Unmarshal(cmd.Payload, &req); err != nil {
				return &ApplyResponse{
					Error: err,
				}
			}
			res, err := b.store.ListStreams(&req)
			return &ApplyResponse{
				Data:  res,
				Error: err,
			}

		case "SIZE":
			var req api.SizeRequest
			if err := proto.Unmarshal(cmd.Payload, &req); err != nil {
				return &ApplyResponse{
					Error: err,
				}
			}
			res, err := b.store.Size(&req)
			return &ApplyResponse{
				Data:  res,
				Error: err,
			}
		default:
			log.Println(cmd.Op)
		}
	}

	return &ApplyResponse{
		Error: errors.New("unknown command"),
	}
}

func (b *fsm) Snapshot() (raft.FSMSnapshot, error) {
	log.Println("SNAPSHOT")
	return newSnapshotNoop()
}

func (b *fsm) Restore(io.ReadCloser) error {
	return nil
}
func NewFSM(store store.EventStore) *fsm {
	return &fsm{store}
}
