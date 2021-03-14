package projections

import (
	"go.etcd.io/bbolt"
)

type ProjectionEngine struct {
	db *bbolt.DB
}

func NewProjectionEngine(db *bbolt.DB) (*ProjectionEngine, error) {
	return &ProjectionEngine{db}, nil
}
