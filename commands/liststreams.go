package commands

import (
	"encoding/json"

	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/kajjagtenberg/go-commando"
)

const (
	CMD_LIST_STREAMS       = "liststreams"
	CMD_LIST_STREAMS_SHORT = "lstreams"
)

type ListStreamsRequest struct {
	Skip  uint32 `json:"skip"`
	Limit uint32 `json:"limit"`
}

type ListStreamsResponse struct {
	Streams []store.Stream `json:"streams"`
}

func ListStreamsHandler(store store.EventStore) commando.CommandHandler {
	return func(cmd commando.Command) (interface{}, error) {
		var req ListStreamsRequest

		if len(cmd.Args) > 0 {
			if err := json.Unmarshal(cmd.Args, &req); err != nil {
				return nil, err
			}
		}

		streams, err := store.ListStreams(req.Skip, req.Limit)
		if err != nil {
			return nil, err
		}

		return ListStreamsResponse{streams}, nil
	}
}

func SetupListStreamsHandler(dispatcher *commando.CommandDispatcher, eventstore store.EventStore) {
	dispatcher.Register(CMD_LIST_STREAMS, CMD_LIST_STREAMS_SHORT, ListStreamsHandler(eventstore))
}
