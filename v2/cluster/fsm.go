package cluster

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	"github.com/kajjagtenberg/eventflowdb/persistence"
	"github.com/oklog/ulid"
)

const (
	BUCKET_STREAMS = "streams"
	BUCKET_EVENTS  = "events"
)

var (
	entropy = ulid.Monotonic(rand.New(rand.NewSource(int64(ulid.Now()))), 0)
)

type ApplyResult struct {
	Value interface{}
	Error error
}

type FSM struct {
	persistence *persistence.Persistence
}

func (fsm *FSM) Apply(applyLog *raft.Log) interface{} {
	var result ApplyResult

	switch applyLog.Type {
	case raft.LogCommand:
		var cmd ApplyLog

		if err := proto.Unmarshal(applyLog.Data, &cmd); err != nil {
			log.Printf("Failed to unmarshal ApplyLog: %v", err)

			result.Error = err

			return result
		}

		switch cmd := cmd.Command.(type) {
		case *ApplyLog_Add:
			var streamID uuid.UUID
			if err := streamID.UnmarshalBinary(cmd.Add.Stream); err != nil {
				result.Error = err

				return result
			}

			version := cmd.Add.Version

			var events []persistence.EventData

			for _, event := range cmd.Add.Events {
				var causationID ulid.ULID
				if event.CausationId != nil {
					if err := causationID.UnmarshalBinary(event.CausationId); err != nil {
						return err
					}
				}

				var correlationID ulid.ULID
				if event.CorrelationId != nil {
					if err := correlationID.UnmarshalBinary(event.CorrelationId); err != nil {
						return err
					}
				}

				events = append(events, persistence.EventData{
					Type:          event.Type,
					Data:          event.Data,
					Metadata:      event.Metadata,
					CausationID:   causationID,
					CorrelationID: correlationID,
					AddedAt:       time.Unix(0, event.AddedAt),
				})
			}

			records, err := fsm.persistence.Add(streamID, version, events)
			if err != nil {
				result.Error = err

				return result
			}

			result.Value = records
		}
	default:
		log.Println("Type:", applyLog.Type)
	}

	return result
}

func (fsm *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, errors.New("Not implemented")
}

func (fsm *FSM) Restore(io.ReadCloser) error {
	return errors.New("Not implemented")
}

func NewFSM(persistence *persistence.Persistence) (*FSM, error) {
	return &FSM{persistence}, nil
}
