package commands

import (
	"github.com/kajjagtenberg/eventflowdb/si"
	"github.com/kajjagtenberg/eventflowdb/store"
)

type SizeResponse struct {
	Size  int64  `json:"size"`
	Human string `json:"human"`
}

func SizeHandler(store store.Store) CommandHandler {
	return func(cmd Command) (interface{}, error) {
		size, err := store.Size()
		if err != nil {
			return nil, err
		}

		return SizeResponse{size, si.ByteCountSI(size)}, nil
	}
}
