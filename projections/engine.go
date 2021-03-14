package projections

import (
	"errors"
	"time"

	"eventflowdb/compiler"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

type ProjectionEngine struct {
	db       *bbolt.DB
	compiler *compiler.Compiler
}

func (engine *ProjectionEngine) CreateProjection(name string, code string) (*Projection, error) {
	if len(name) == 0 {
		return nil, errors.New("Name cannot be empty")
	}

	if len(code) == 0 {
		return nil, errors.New("Code cannot be empty")
	}

	compiled, err := engine.compiler.Compile(code)
	if err != nil {
		return nil, err
	}

	proj := &Projection{
		ID:           uuid.New(),
		Name:         name,
		Code:         code,
		CompiledCode: compiled,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := engine.db.Update(func(t *bbolt.Tx) error {
		return nil
	}); err != nil {
		return nil, err
	}

	return proj, nil
}

func (engine *ProjectionEngine) GetProjections() ([]Projection, error) {
	var projections []Projection

	if err := engine.db.View(func(t *bbolt.Tx) error {

		return nil
	}); err != nil {
		return nil, err
	}

	return projections, nil
}

func NewProjectionEngine(db *bbolt.DB) (*ProjectionEngine, error) {
	if err := db.Update(func(t *bbolt.Tx) error {
		if _, err := t.CreateBucketIfNotExists([]byte("projections")); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	compiler, err := compiler.NewCompiler()
	if err != nil {
		return nil, err
	}

	return &ProjectionEngine{db, compiler}, nil
}
