package store

import "errors"

var (
	ErrConcurrentStreamModification = errors.New("concurrent stream modification")
	ErrGappedStream                 = errors.New("given version leaves gap in stream")
	ErrEmptyEventType               = errors.New("event type cannot be empty")
)
