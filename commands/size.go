package commands

import (
	"github.com/kajjagtenberg/eventflowdb/si"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/kajjagtenberg/go-commando"
)

const (
	CMD_SIZE       = "size"
	CMD_SIZE_SHORT = "s"
)

type SizeResponse struct {
	Size  int64  `json:"size"`
	Human string `json:"human"`
}

func SizeHandler(store store.EventStore) commando.CommandHandler {
	return func(cmd commando.Command) (interface{}, error) {
		size, err := store.Size()
		if err != nil {
			return nil, err
		}

		return SizeResponse{size, si.ByteCountSI(size)}, nil
	}
}
