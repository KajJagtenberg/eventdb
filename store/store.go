package store

import (
	"io"

	"github.com/eventflowdb/eventflowdb/api"
)

type EventStore interface {
	GetStream(*api.GetStreamRequest) (*api.GetStreamResponse, error)

	GetGlobalStream(*api.GetGlobalStreamRequest) (*api.GetGlobalStreamResponse, error)

	AppendStream(*api.AppendStreamRequest) (*api.AppendStreamResponse, error)

	GetEvent(*api.GetEventRequest) (*api.Event, error)

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
	EventCountEstimate(*api.EventCountEstimateRequest) (*api.EventCountResponse, error)

	/*
		Returns an estimate of the total number of streams in the database
	*/
	StreamCountEstimate(*api.StreamCountEstimateRequest) (*api.StreamCountResponse, error)

	/*
		Returns the names of streams, with skip and limit options
	*/

	ListStreams(*api.ListStreamsRequest) (*api.ListStreamsReponse, error)

	// Checksum(*api.ChecksumRequest) (*api.ChecksumResponse, error)

	Close() error
}
