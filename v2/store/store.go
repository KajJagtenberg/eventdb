package store

import (
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	entropy = ulid.Monotonic(rand.New(rand.NewSource(int64(ulid.Now()))), 0)
)

type Storage struct {
	db *bbolt.DB
}

func (s *Storage) Add(req *AddRequest) ([]*RecordedEvent, error) {
	var stream uuid.UUID
	if err := stream.UnmarshalBinary(req.Stream); err != nil {
		err := status.Error(codes.InvalidArgument, "Invalid stream uuid")
		return nil, err
	}

	var events []*RecordedEvent

	for i, event := range req.Events {
		id, err := ulid.New(ulid.Now(), entropy)
		if err != nil {
			return nil, err
		}

		bId, err := id.MarshalBinary()
		if err != nil {
			return nil, err
		}

		if event.CausationId == nil {
			event.CausationId = bId
		}

		if event.CorrelationId == nil {
			event.CorrelationId = bId
		}

		events = append(events, &RecordedEvent{
			Id:            bId,
			Stream:        req.Stream,
			Version:       req.Version + uint32(i),
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   event.CausationId,
			CorrelationId: event.CorrelationId,
			AddedAt:       time.Now().UnixNano(),
		})
	}

	return events, nil
}

func (s *Storage) Get(req *GetRequest) ([]*RecordedEvent, error) {
	var stream uuid.UUID
	if err := stream.UnmarshalBinary(req.Stream); err != nil {
		err := status.Error(codes.InvalidArgument, "Invalid stream uuid")
		return nil, err
	}

	return nil, errors.New("Not implemented")
}

func (s *Storage) Log(req *LogRequest) ([]*RecordedEvent, error) {
	// TODO: Validate offset

	return nil, nil
}

func NewStorage(db *bbolt.DB) *Storage {
	return &Storage{db}
}
