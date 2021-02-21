package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/base64"
	"eventflowdb/graph/generated"
	"eventflowdb/graph/model"
	"eventflowdb/store"
	"log"

	"github.com/google/uuid"
)

func (r *mutationResolver) Append(ctx context.Context, stream string, version int, events []*model.EventData) ([]*model.Event, error) {
	name, err := uuid.Parse(stream)
	if err != nil {
		return nil, err
	}

	var data []store.EventData

	for _, event := range events {
		log.Println(event)

		decodedData, err := base64.StdEncoding.DecodeString(event.Data)
		if err != nil {
			return nil, err
		}

		var metadata []byte

		if event.Metadata != nil {
			metadata, err = base64.StdEncoding.DecodeString(*event.Metadata)
		}

		if err != nil {
			return nil, err
		}

		data = append(data, store.EventData{
			Type:     event.Type,
			Data:     decodedData,
			Metadata: metadata,
		})
	}

	records, err := r.EventStore.AppendToStream(name, version, data)
	if err != nil {
		return nil, err
	}

	var result []*model.Event

	for _, record := range records {
		result = append(result, &model.Event{
			ID:       record.ID.String(),
			Stream:   record.Stream.String(),
			Version:  record.Version,
			Type:     record.Type,
			Data:     base64.StdEncoding.EncodeToString(record.Data),
			Metadata: base64.StdEncoding.EncodeToString(record.Metadata),
			AddedAt:  record.AddedAt,
		})
	}

	return result, nil
}

func (r *queryResolver) FromStream(ctx context.Context, stream *string, version *int, limit *int) ([]*model.Event, error) {
	return nil, nil
}

func (r *queryResolver) FromAllStreams(ctx context.Context, offset *string, limit *int) ([]*model.Event, error) {
	return nil, nil
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
func (r *queryResolver) LoadFromStream(ctx context.Context, stream *string, version *int, limit *int) ([]*model.Event, error) {
	return nil, nil
}
