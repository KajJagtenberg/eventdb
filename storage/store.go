package storage

import (
	"log"
	"time"

	"github.com/eventflowdb/eventflowdb/api"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SQLEventStore struct {
	db *gorm.DB
}

func (s *SQLEventStore) GetStream(req *api.GetStreamRequest) (*api.GetStreamResponse, error) {
	res := &api.GetStreamResponse{}

	stream, err := uuid.Parse(req.Stream)
	if err != nil {
		return nil, err
	}

	var events []Event

	if err := s.db.Where("stream = ? AND version >= ?", stream, req.Version).Limit(int(req.Limit)).Select("id").Find(&events).Error; err != nil {
		return nil, err
	}

	for _, event := range events {
		res.Events = append(res.Events, event.ID.String())
	}

	return res, nil
}

func (s *SQLEventStore) GetGlobalStream(req *api.GetGlobalStreamRequest) (*api.GetGlobalStreamResponse, error) {
	res := &api.GetGlobalStreamResponse{}

	var events []Event

	if err := s.db.Where("timestamp > ?", req.Offset).Limit(int(req.Limit)).Select("id").Find(&events).Error; err != nil {
		return nil, err
	}

	log.Println(events)

	for _, event := range events {
		res.Events = append(res.Events, event.ID.String())
	}

	return res, nil
}

func (s *SQLEventStore) AppendToStream(req *api.AppendToStreamRequest) (*api.AppendToStreamResponse, error) {
	res := &api.AppendToStreamResponse{}

	if len(req.Events) == 0 {
		return nil, ErrEmptyEvents
	}

	stream, err := uuid.Parse(req.Stream)
	if err != nil {
		return nil, err
	}

	var version int

	if err := s.db.Raw("SELECT COALESCE(MAX(version), -1) AS version FROM events WHERE stream = ?;", stream).Scan(&version).Error; err != nil {
		return nil, err
	}

	if int(req.Version) != version+1 {
		return nil, ErrConcurrentStreamModification
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		for i, event := range req.Events {
			id := uuid.New()

			var causationID uuid.UUID

			if len(event.CausationId) > 0 {
				var err error
				causationID, err = uuid.Parse(event.CausationId)

				if err != nil {
					return err
				}
			} else {
				causationID = id
			}

			var correlationID uuid.UUID

			if len(event.CorrelationId) > 0 {
				var err error
				correlationID, err = uuid.Parse(event.CorrelationId)

				if err != nil {
					return err
				}
			} else {
				correlationID = id
			}

			if err := tx.Create(&Event{
				ID:            id,
				Stream:        stream,
				Version:       uint32(version) + 1 + uint32(i),
				Type:          event.Type,
				Data:          event.Data,
				Metadata:      event.Metadata,
				CausationID:   causationID,
				CorrelationID: correlationID,
				Timestamp:     time.Now().Unix(),
			}).Error; err != nil {
				return err
			}

			res.Events = append(res.Events, id.String())
		}

		return nil
	}); err != nil {
		return nil, err
	}
	// TODO: handle duplicate error

	return res, nil
}

func (s *SQLEventStore) GetEvent(req *api.GetEventRequest) (*api.Event, error) {
	var event Event

	if err := s.db.Where("id = ?", req.Id).First(&event).Error; err != nil {
		return nil, err
	}

	return &api.Event{
		Id:            event.ID.String(),
		Stream:        event.Stream.String(),
		Version:       event.Version,
		Type:          event.Type,
		Data:          event.Data,
		Metadata:      event.Metadata,
		CausationId:   event.CausationID.String(),
		CorrelationId: event.CorrelationID.String(),
		AddedAt:       event.Timestamp,
	}, nil
}

func (s *SQLEventStore) Size(*api.SizeRequest) (*api.SizeResponse, error) {
	return nil, nil
}

func (s *SQLEventStore) EventCount(*api.EventCountRequest) (*api.EventCountResponse, error) {
	var count int64

	if err := s.db.Raw("SELECT COUNT(id) FROM events;").First(&count).Error; err != nil {
		return nil, err
	}

	return &api.EventCountResponse{
		Count: count,
	}, nil
}

func (s *SQLEventStore) StreamCount(*api.StreamCountRequest) (*api.StreamCountResponse, error) {
	return nil, nil
}

func (s *SQLEventStore) ListStreams(*api.ListStreamsRequest) (*api.ListStreamsReponse, error) {
	return nil, nil
}

func (s *SQLEventStore) Close() error {
	return nil
}

func NewSQLEventStore(db *gorm.DB) (*SQLEventStore, error) {
	if err := db.AutoMigrate(&Event{}); err != nil {
		return nil, err
	}

	return &SQLEventStore{db}, nil
}

type Event struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Stream        uuid.UUID `gorm:"type:uuid;index:idx_stream,unique" json:"stream"`
	Version       uint32    `gorm:"index:idx_stream,unique" json:"version"`
	Type          string    `json:"type"`
	Data          []byte    `json:"data"`
	Metadata      []byte    `json:"metadata"`
	CausationID   uuid.UUID `json:"causation_id"`
	CorrelationID uuid.UUID `json:"correlation_id"`
	Timestamp     int64     `json:"timestamp"`
}
