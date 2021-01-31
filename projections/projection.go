package projections

import "github.com/oklog/ulid"

type Projection struct {
}

type ProjectionConfig struct {
	Checkpoint   ulid.ULID
	Code         string
	CompiledCode string
}
