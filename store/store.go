package store

import (
	"io"

	"github.com/kajjagtenberg/eventflowdb/api"
)

type EventStore interface {
	/*
		Adds event to specified stream at specfied version offset. Returns the persisted events and
		an error in case of concurrent stream modification.
	*/
	// Add(stream uuid.UUID, version uint32, events []EventData) ([]Event, error)
	Add(*api.AddRequest) (*api.EventResponse, error)

	/*
		Returns events for specified stream, offset at given version and limits the resulting set by the given limit.
		If limit is zero, then all events from given version onwards will be returned.
	*/
	Get(*api.GetRequest) (*api.EventResponse, error)

	/*
		Returns amount of events that have been recorded since the given offset. The maximum amount of returned events is specified by the given limit.
		If the limit is 0, then it will return a maximum of 100 events.
	*/
	GetAll(*api.GetAllRequest) (*api.EventResponse, error)

	/*
		Size of the database in bytes on disk
	*/
	Size(*api.SizeRequest) (*api.SizeResponse, error)

	/*
		Writes a snapshot of the database to a writer
	*/
	Backup(dst io.Writer) error

	/*
		Returns the total number of events stored in the database
	*/
	EventCount(*api.EventCountRequest) (*api.EventCountResponse, error)

	/*
		Returns the total number of streams in the database
	*/
	StreamCount(*api.StreamCountRequest) (*api.StreamCountResponse, error)

	/*
		Returns an estimate of the total number of events stored in the database
	*/
	EventCountEstimate(*api.EventCountRequest) (*api.EventCountResponse, error)

	/*
		Returns an estimate of the total number of streams in the database
	*/
	StreamCountEstimate(*api.StreamCountRequest) (*api.StreamCountResponse, error)

	/*
		Returns the names of streams, with skip and limit options
	*/

	ListStreams(*api.ListStreamsRequest) (*api.ListStreamsReponse, error)

	Checksum(*api.ChecksumRequest) (*api.ChecksumResponse, error)

	Close() error
}
