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

	/*
		Returns events for specified stream, offset at given version and limits the resulting set by the given limit.
		If limit is zero, then all events from given version onwards will be returned.
	*/
	Get(stream uuid.UUID, version uint32, limit uint32) ([]Event, error)
}
