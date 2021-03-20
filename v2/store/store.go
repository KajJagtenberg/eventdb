package store

import (
	"bytes"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	entropy = ulid.Monotonic(rand.New(rand.NewSource(int64(ulid.Now()))), 0)
)

var (
	ErrConcurrentStreamModification = errors.New("Concurrent stream modification")
	ErrNoEvents                     = errors.New("No events specified")
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

	if len(req.Events) == 0 {
		return nil, ErrNoEvents
	}

	var records []*RecordedEvent

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

		records = append(records, &RecordedEvent{
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

	if err := s.db.Batch(func(t *bbolt.Tx) error {
		streams := t.Bucket([]byte("streams"))
		events := t.Bucket([]byte("events"))

		stream := &Stream{}

		v := streams.Get(req.Stream)

		if v == nil {
			stream.Id = req.Stream
			stream.AddedAt = time.Now().UnixNano()
		} else {
			if err := proto.Unmarshal(v, stream); err != nil {
				log.Printf("Unable to decode stream: %v", err)
				return err
			}
		}

		if len(stream.Events) != int(req.Version) {
			return ErrConcurrentStreamModification
		}

		for _, record := range records {
			v, err := proto.Marshal(record)
			if err != nil {
				log.Printf("Unable to encode event: %v", err)
				return err
			}

			if err := events.Put(record.Id, v); err != nil {
				return err
			}

			stream.Events = append(stream.Events, record.Id)
		}

		v, err := proto.Marshal(stream)
		if err != nil {
			return err
		}

		if err := streams.Put(req.Stream, v); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return records, nil
}

func (s *Storage) Get(req *GetRequest) ([]*RecordedEvent, error) {
	var stream uuid.UUID
	if err := stream.UnmarshalBinary(req.Stream); err != nil {
		err := status.Error(codes.InvalidArgument, "Invalid stream uuid")
		return nil, err
	}

	var result []*RecordedEvent

	if err := s.db.View(func(t *bbolt.Tx) error {
		streams := t.Bucket([]byte("streams"))
		events := t.Bucket([]byte("events"))

		v := streams.Get(req.Stream)

		if v == nil {
			return nil
		}

		stream := &Stream{}

		if err := proto.Unmarshal(v, stream); err != nil {
			log.Printf("Unable to decode stream: %v", err)
			return nil
		}

		for _, id := range stream.Events {
			if req.Version > 0 {
				req.Version--
				continue
			}
			if req.Limit != 0 && len(result) >= int(req.Limit) {
				break
			}

			v := events.Get(id)

			event := &RecordedEvent{}

			if err := proto.Unmarshal(v, event); err != nil {
				log.Printf("Unable to decode event: %v", err)
				return err
			}

			result = append(result, event)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) Log(req *LogRequest) ([]*RecordedEvent, error) {
	var offset ulid.ULID
	if err := offset.UnmarshalBinary(req.Offset); err != nil {
		return nil, errors.New("Unable to decode offset to a valid ULID")
	}

	var result []*RecordedEvent

	if err := s.db.View(func(t *bbolt.Tx) error {
		events := t.Bucket([]byte("events")).Cursor()

		for k, v := events.Seek(req.Offset); k != nil && (req.Limit == 0 || len(result) >= int(req.Limit)); k, v = events.Next() {
			if bytes.Compare(k, req.Offset) == 0 {
				continue
			}

			event := &RecordedEvent{}

			if err := proto.Unmarshal(v, event); err != nil {
				log.Printf("Unable to decode event: %v", err)
				return err
			}

			result = append(result, event)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func NewStorage(db *bbolt.DB) (*Storage, error) {
	if err := db.Update(func(t *bbolt.Tx) error {
		buckets := []string{"events", "streams"}

		for _, bucket := range buckets {
			if _, err := t.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return &Storage{db}, nil
}
