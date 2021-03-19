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
	streams map[uuid.UUID][]*RecordedEvent
	log     map[ulid.ULID]*RecordedEvent
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

	s.streams[stream] = append(s.streams[stream], events...)

	return &AddResponse{
		Events: events,
	}, nil
}

func (s *StoreService) Get(in *GetRequest, result EventStore_GetServer) error {
	var stream uuid.UUID
	if err := stream.UnmarshalBinary(in.Stream); err != nil {
		err := status.Error(codes.InvalidArgument, "Invalid stream uuid")
		return err
	}

	events := s.streams[stream]

	sent := 0

	for _, event := range events {
		if in.Version > 0 {
			in.Version--
			continue
		}

		if sent >= int(in.Limit) {
			break
		}

		if err := result.Send(event); err != nil {
			return err
		}

		sent++
	}

	return nil
}

func (s *StoreService) Log(in *LogRequest, result EventStore_LogServer) error {
	return nil
}

func NewStoreService() *StoreService {
	return &StoreService{
		entropy: ulid.Monotonic(rand.New(rand.NewSource(int64(ulid.Now()))), 0),
		streams: make(map[uuid.UUID][]*RecordedEvent),
		log:     make(map[ulid.ULID]*RecordedEvent),
	}
}
