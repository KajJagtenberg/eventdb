package store

import (
	"log"
	"time"

	"github.com/eventflowdb/eventflowdb/api"
	"github.com/google/uuid"
)

func convertToEvents(req *api.AppendToStreamRequest) (result []*api.Event, err error) {
	for i, event := range req.Events {
		if len(event.Type) == 0 {
			return nil, ErrEmptyEventType
		}

		id := uuid.New()

		if len(event.CausationId) == 0 {
			event.CausationId = id.String()
		}
		if len(event.CorrelationId) == 0 {
			event.CorrelationId = id.String()
		}

		result = append(result, &api.Event{
			Id:            id.String(),
			Stream:        req.Stream,
			Version:       req.Version + int32(i),
			Type:          event.Type,
			Metadata:      event.Metadata,
			CausationId:   event.CausationId,
			CorrelationId: event.CorrelationId,
			Data:          event.Data,
			AddedAt:       time.Now().Unix(),
		})
	}
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}
