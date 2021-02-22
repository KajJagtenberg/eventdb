package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"eventflowdb/graph/generated"
	"eventflowdb/graph/model"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"go.etcd.io/bbolt"
)

func (r *mutationResolver) Append(ctx context.Context, stream string, version int, events []*model.EventData) ([]*model.Event, error) {
	name, err := uuid.Parse(stream)
	if err != nil {
		return nil, err
	}

	if version < 0 {
		return nil, errors.New("Version cannot be negative")
	}

	if len(events) == 0 {
		return nil, errors.New("List of event data cannot be empty")
	}

	result := []*model.Event{}

	err = r.DB.Update(func(txn *bbolt.Tx) error {
		streamBucket := txn.Bucket([]byte("streams"))
		eventBucket := txn.Bucket([]byte("events"))

		key := name[:]
		value := streamBucket.Get(key)

		var stream []ulid.ULID

		if value != nil {
			if err := json.Unmarshal(value, &stream); err != nil {
				return err
			}
		}

		if len(stream) != version {
			return ErrConcurrentStreamModification
		}

		for i, eventData := range events {
			id, err := ulid.New(ulid.Now(), entropy)
			if err != nil {
				return err
			}

			stream = append(stream, id)

			event := &model.Event{
				ID:       id.String(),
				Stream:   name.String(),
				Version:  version + i,
				Type:     eventData.Type,
				Data:     eventData.Data,
				Metadata: eventData.Metadata,
				AddedAt:  time.Now(),
			}

			result = append(result, event)

			raw, err := json.Marshal(event)
			if err != nil {
				return err
			}

			if err := eventBucket.Put(id[:], raw); err != nil {
				return err
			}
		}

		value, err = json.Marshal(stream)
		if err != nil {
			return err
		}

		if err := streamBucket.Put(key, value); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) Streams(ctx context.Context, skip int, limit int) ([]*model.Stream, error) {
	streams := []*model.Stream{}

	err := r.DB.View(func(txn *bbolt.Tx) error {
		streamBucket := txn.Bucket([]byte("streams"))
		cursor := streamBucket.Cursor()

		for k, v := cursor.First(); k != nil && (len(streams) < limit || limit == 0); k, v = cursor.Next() {
			if skip > 0 {
				skip--
				continue
			}

			name, err := uuid.FromBytes(k)
			if err != nil {
				return err
			}

			stream := []ulid.ULID{}

			if err := json.Unmarshal(v, &stream); err != nil {
				return err
			}

			streams = append(streams, &model.Stream{
				Name: name.String(),
				Size: len(stream),
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return streams, nil
}

func (r *queryResolver) Stream(ctx context.Context, id string, skip int, limit int) ([]*model.Event, error) {
	name, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	if skip < 0 {
		return nil, errors.New("Skip cannot be negative")
	}

	if limit < 0 {
		return nil, errors.New("Limit cannot be negative")
	}

	events := []*model.Event{}

	err = r.DB.View(func(txn *bbolt.Tx) error {
		streamBucket := txn.Bucket([]byte("streams"))
		eventsBucket := txn.Bucket([]byte("events"))

		key := name[:]
		value := streamBucket.Get(key)

		stream := []ulid.ULID{}

		if err := json.Unmarshal(value, &stream); err != nil {
			return err
		}

		for _, id := range stream {
			if skip > 0 {
				skip--
				continue
			}

			if len(events) == limit && limit != 0 {
				break
			}

			key := id[:]
			value := eventsBucket.Get(key)

			event := &model.Event{}

			if err := json.Unmarshal(value, &event); err != nil {
				return err
			}

			events = append(events, event)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *queryResolver) All(ctx context.Context, offset string, limit int) ([]*model.Event, error) {
	var name ulid.ULID
	var err error

	if len(offset) > 0 {
		name, err = ulid.Parse(offset)
		if err != nil {
			return nil, err
		}
	}

	if limit < 0 {
		return nil, errors.New("Limit cannot be negative")
	}

	events := []*model.Event{}

	err = r.DB.View(func(txn *bbolt.Tx) error {
		eventsBucket := txn.Bucket([]byte("events"))
		cursor := eventsBucket.Cursor()

		for k, v := cursor.Seek(name[:]); k != nil && (len(events) < limit || limit == 0); k, v = cursor.Next() {
			if bytes.Compare(k, name[:]) == 0 {
				continue
			}

			event := &model.Event{}

			if err := json.Unmarshal(v, &event); err != nil {
				return err
			}

			events = append(events, event)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *queryResolver) TotalStreams(ctx context.Context) (int, error) {
	total := 0

	err := r.DB.View(func(txn *bbolt.Tx) error {
		cursor := txn.Bucket([]byte("streams")).Cursor()

		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			total++
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *queryResolver) TotalEvents(ctx context.Context) (int, error) {
	total := 0

	err := r.DB.View(func(txn *bbolt.Tx) error {
		cursor := txn.Bucket([]byte("events")).Cursor()

		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			total++
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
var (
	ErrNotImplemented               = errors.New("Not implemented")
	ErrConcurrentStreamModification = errors.New("Concurrent stream modification")

	entropy = ulid.Monotonic(rand.New(rand.NewSource(int64(ulid.Now()))), 0)
)
