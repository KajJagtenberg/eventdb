package store

import (
	"io"

	"github.com/google/uuid"
)

type Store interface {
	/*
		Size of the database in bytes on disk
	*/
	Size() int64

	/*
		Writes a snapshot of the database to a writer
	*/
	Backup(dst io.Writer) error

	/*
		Adds event to specified stream at specfied version offset. Returns the persisted events and
		an error in case of concurrent stream modification.
	*/
	Add(stream uuid.UUID, version uint32, events []EventData) ([]Event, error)
}
