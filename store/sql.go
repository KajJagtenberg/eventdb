package store

import (
	"time"

	"github.com/eventflowdb/eventflowdb/api"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SQLStore struct {
	db *gorm.DB
}

func (s *SQLStore) GetStream(*api.GetStreamRequest) (*api.GetStreamResponse, error) {
	return nil, nil
}

func (s *SQLStore) GetGlobalStream(*api.GetGlobalStreamRequest) (*api.GetGlobalStreamResponse, error) {
	return nil, nil
}

func (s *SQLStore) AppendToStream(req *api.AppendToStreamRequest) (*api.AppendToStreamResponse, error) {
	res := &api.AppendToStreamResponse{}

	if len(req.Events) == 0 {
		return nil, ErrEmptyEvents
	}

	stream, err := uuid.Parse(req.Stream)
	if err != nil {
		return nil, err
	}

	var version int32

	if err := s.db.Raw("SELECT COALESCE(MAX(version), -1) AS version FROM events WHERE stream = ?;", stream).Scan(&version).Error; err != nil {
		return nil, err
	}

	if req.Version != version+1 {
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
				Version:       version + 1 + int32(i),
				Type:          event.Type,
				Data:          event.Data,
				Metadata:      event.Metadata,
				CausationID:   causationID,
				CorrelationID: correlationID,
				Timestamp:     time.Now().Unix(),
			}).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}
	// TODO: handle duplicate error

	return res, nil
}

func (s *SQLStore) GetEvent(*api.GetEventRequest) (*api.Event, error) {
	return nil, nil
}

func (s *SQLStore) Size(*api.SizeRequest) (*api.SizeResponse, error) {
	return nil, nil
}

func (s *SQLStore) EventCount(*api.EventCountRequest) (*api.EventCountResponse, error) {
	return nil, nil
}

func (s *SQLStore) StreamCount(*api.StreamCountRequest) (*api.StreamCountResponse, error) {
	return nil, nil
}

func (s *SQLStore) ListStreams(*api.ListStreamsRequest) (*api.ListStreamsReponse, error) {
	return nil, nil
}

func (s *SQLStore) Close() error {
	return nil
}

func NewSQLStore(db *gorm.DB) (*SQLStore, error) {
	if err := db.AutoMigrate(&Event{}); err != nil {
		return nil, err
	}

	return &SQLStore{db}, nil
}

type Event struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Stream        uuid.UUID `gorm:"type:uuid;index:idx_stream,unique"`
	Version       int32     `gorm:"index:idx_stream,unique"`
	Type          string
	Data          []byte
	Metadata      []byte
	CausationID   uuid.UUID
	CorrelationID uuid.UUID
	Timestamp     int64
}
