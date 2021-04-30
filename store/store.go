package store

import (
	"io"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

type Store interface {
	/*
		Size of the database in bytes on disk
	*/
	Size() (int64, error)

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

	/*
		Returns amount of events that have been recorded since the given offset. The maximum amount of returned events is specified by the given limit.
		If the limit is 0, then it will return a maximum of 100 events.
	*/
	GetAll(offset ulid.ULID, limit uint32) ([]Event, error)

	/*
		Returns the total number of events stored in the database
	*/
	EventCount() (int64, error)

	/*
		Returns the total number of streams in the database
	*/
	StreamCount() (int64, error)

	/*
		Returns an estimate of the total number of events stored in the database
	*/
	EventCountEstimate() (int64, error)

	/*
		Returns an estimate of the total number of streams in the database
	*/
	StreamCountEstimate() (int64, error)

	Close() error
}
