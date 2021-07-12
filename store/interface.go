package store

import (
	"github.com/eventflowdb/eventflowdb/api"
)

type EventStore interface {
	GetStream(*api.GetStreamRequest) (*api.GetStreamResponse, error)

	GetGlobalStream(*api.GetGlobalStreamRequest) (*api.GetGlobalStreamResponse, error)

	AppendToStream(*api.AppendToStreamRequest) (*api.AppendToStreamResponse, error)

	GetEvent(*api.GetEventRequest) (*api.Event, error)

	Size(*api.SizeRequest) (*api.SizeResponse, error)

	EventCount(*api.EventCountRequest) (*api.EventCountResponse, error)

	StreamCount(*api.StreamCountRequest) (*api.StreamCountResponse, error)

	ListStreams(*api.ListStreamsRequest) (*api.ListStreamsReponse, error)

	Close() error
}
