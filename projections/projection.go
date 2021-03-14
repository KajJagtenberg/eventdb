package projections

import "github.com/google/uuid"

type Projection struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	CompiledCode string    `json:"compiled_code"`
}
