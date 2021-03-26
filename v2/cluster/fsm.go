package cluster

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/raft"
	"github.com/oklog/ulid"
	"go.etcd.io/bbolt"
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
	db *bbolt.DB
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
			events, err := fsm.add(cmd.Add)

			result.Value = events
			result.Error = err
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

func (fsm *FSM) add(cmd *AddCommand) ([]Event, error) {
	var events []Event

	err := fsm.db.Update(func(t *bbolt.Tx) error {
		streams := t.Bucket([]byte(BUCKET_STREAMS))
		events := t.Bucket([]byte(BUCKET_EVENTS))

		var stream Stream

		if packed := streams.Get(cmd.Stream); packed != nil {
			if err := proto.Unmarshal(packed, &stream); err != nil {
				return err
			}
		} else {
			stream.AddedAt = time.Now().UnixNano()
			stream.Id = cmd.Stream
		}

		if len(stream.Events) != int(cmd.Version) {
			return errors.New("Concurrent stream modification")
		}

		for i, event := range cmd.Events {
			id, err := ulid.New(ulid.Now(), entropy)
			if err != nil {
				return err
			}

			record := Event{}
			record.Id = id[:]
			record.Stream = cmd.Stream
			record.Version = cmd.Version + uint32(i)
			record.Data = event.Data
			record.Metadata = event.Metadata
			record.CausationId = event.CausationId
			record.CorrelationId = event.CorrelationId

			if record.CausationId == nil {
				record.CausationId = record.Id
			}

			if record.CorrelationId == nil {
				record.CorrelationId = record.Id
			}

			record.AddedAt = time.Now().UnixNano()

			stream.Events = append(stream.Events, record.Id)

			packed, err := proto.Marshal(&record)
			if err != nil {
				return err
			}

			if err := events.Put(record.Id[:], packed); err != nil {
				return err
			}
		}

		packed, err := proto.Marshal(&stream)
		if err != nil {
			return err
		}

		return streams.Put(stream.Id[:], packed)
	})

	if err != nil {
		return nil, err
	}

	return events, nil
}

func NewFSM(db *bbolt.DB) (*FSM, error) {
	err := db.Update(func(t *bbolt.Tx) error {
		buckets := []string{BUCKET_STREAMS, BUCKET_EVENTS}

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
