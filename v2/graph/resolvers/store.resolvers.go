package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/graph/model"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/oklog/ulid"
)

func (r *queryResolver) StreamCount(ctx context.Context) (int, error) {
	streamCount, err := r.Storage.StreamCount()
	if err != nil {
		return 0, err
	}

	return int(streamCount), nil
}

func (r *queryResolver) Get(ctx context.Context, input model.GetInput) ([]*model.RecordedEvent, error) {
	stream, err := uuid.Parse(input.Stream)
	if err != nil {
		return nil, err
	}

	records, err := r.Storage.Get(&store.GetRequest{
		Stream: stream[:],
	})
	if err != nil {
		return nil, err
	}

	result := []*model.RecordedEvent{}

	for _, record := range records {
		var id ulid.ULID
		if err := id.UnmarshalBinary(record.Id); err != nil {
			return nil, err
		}

		var causation_id ulid.ULID
		if err := causation_id.UnmarshalBinary(record.CausationId); err != nil {
			return nil, err
		}

		var correlation_id ulid.ULID
		if err := correlation_id.UnmarshalBinary(record.CorrelationId); err != nil {
			return nil, err
		}

		result = append(result, &model.RecordedEvent{
			ID:            id.String(),
			Stream:        stream.String(),
			Version:       int(record.Version),
			Data:          base64.StdEncoding.EncodeToString(record.Data),
			Metadata:      base64.StdEncoding.EncodeToString(record.Metadata),
			CausationID:   causation_id.String(),
			CorrelationID: correlation_id.String(),
			AddedAt:       time.Unix(0, record.AddedAt),
		})
	}

	return result, nil
}

func (r *queryResolver) Log(ctx context.Context, input model.LogInput) ([]*model.RecordedEvent, error) {
	var offset ulid.ULID
	var err error

	if len(input.Offset) > 0 {
		offset, err = ulid.Parse(input.Offset)
		if err != nil {
			return nil, err
		}
	}
	records, err := r.Storage.Log(&store.LogRequest{
		Offset: offset[:],
		Limit:  uint32(input.Limit),
	})
	if err != nil {
		return nil, err
	}

	result := []*model.RecordedEvent{}

	for _, record := range records {
		var id ulid.ULID
		if err := id.UnmarshalBinary(record.Id); err != nil {
			return nil, err
		}

		var causation_id ulid.ULID
		if err := causation_id.UnmarshalBinary(record.CausationId); err != nil {
			return nil, err
		}

		var correlation_id ulid.ULID
		if err := correlation_id.UnmarshalBinary(record.CorrelationId); err != nil {
			return nil, err
		}

		var stream uuid.UUID
		if err := stream.UnmarshalBinary(record.Stream); err != nil {
			return nil, err
		}

		result = append(result, &model.RecordedEvent{
			ID:            id.String(),
			Stream:        stream.String(),
			Version:       int(record.Version),
			Data:          base64.StdEncoding.EncodeToString(record.Data),
			Metadata:      base64.StdEncoding.EncodeToString(record.Metadata),
			CausationID:   causation_id.String(),
			CorrelationID: correlation_id.String(),
			AddedAt:       time.Unix(0, record.AddedAt),
		})
	}

	return result, nil
}

func (r *queryResolver) Streams(ctx context.Context, input model.StreamsInput) ([]*model.Stream, error) {
	streams, err := r.Storage.GetStreams(uint32(input.Skip), uint32(input.Limit))
	if err != nil {
		return nil, err
	}

	result := []*model.Stream{}

	for _, stream := range streams {
		id, err := uuid.FromBytes(stream.Id)
		if err != nil {
			return nil, err
		}

		result = append(result, &model.Stream{
			ID: id.String(),
		})
	}

	return result, nil
}
