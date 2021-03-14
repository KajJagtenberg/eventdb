package projections

import (
	"go.etcd.io/bbolt"
)

type ProjectionEngine struct {
	db *bbolt.DB
}

func (engine *ProjectionEngine) GetProjections() ([]Projection, error) {
	return []Projection{}, nil
}

func NewProjectionEngine(db *bbolt.DB) (*ProjectionEngine, error) {
	return &ProjectionEngine{db}, nil
}
