package projections

import (
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

type Projection struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	CompiledCode string    `json:"compiled_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Checkpoint   ulid.ULID `json:"checkpoint"`
}
