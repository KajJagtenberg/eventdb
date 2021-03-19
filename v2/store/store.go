package store

import (
	"context"
	"io"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StoreService struct {
	entropy io.Reader
}

func (s *StoreService) Add(ctx context.Context, in *AddRequest) (*AddResponse, error) {
	var stream uuid.UUID
	if err := stream.UnmarshalBinary(in.Stream); err != nil {
		err := status.Error(codes.InvalidArgument, "Invalid stream uuid")
		return nil, err
	}

	var events []*RecordedEvent

	for i, event := range in.Events {
		id, err := ulid.New(ulid.Now(), s.entropy)
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
			Stream:        in.Stream,
			Version:       in.Version + uint32(i),
			Type:          event.Type,
			Data:          event.Data,
			Metadata:      event.Metadata,
			CausationId:   event.CausationId,
			CorrelationId: event.CorrelationId,
			AddedAt:       time.Now().UnixNano(),
		})
	}

	return &AddResponse{
		Events: events,
	}, nil
}

func (s *StoreService) Get(ctx context.Context, in *GetRequest) (*GetResponse, error) {
	return nil, nil
}

func NewStoreService() *StoreService {
	return &StoreService{
		entropy: ulid.Monotonic(rand.New(rand.NewSource(int64(ulid.Now()))), 0),
	}
}
