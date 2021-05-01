package api

import (
	"crypto/rand"
	"time"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/store"
)

func Populate(s store.Store, c *Ctx) error {
	count := 100
	size := 100

	for i := 0; i < count; i++ {
		stream := uuid.New()

		data := make([]byte, size)
		if _, err := rand.Read(data); err != nil {
			return err
		}

		events := []store.EventData{
			{
				Type:    "RandomEvent",
				Data:    data,
				AddedAt: time.Now(),
			},
		}

		if _, err := s.Add(stream, 0, events); err != nil {
			return err
		}
	}

	c.Conn.WriteString("OK")

	return nil
}
