package commands

import (
	"encoding/base32"

	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/oklog/ulid"
)

const (
	CMD_CHECKSUM       = "checksum"
	CMD_CHECKSUM_SHORT = "ch"
)

type ChecksumResponse struct {
	ID       ulid.ULID `json:"id"`
	Checksum string    `json:"checksum"`
}

func ChecksumHandler(s store.Store) CommandHandler {
	return func(cmd Command) (interface{}, error) {
		id, checksum, err := s.Checksum()
		if err != nil {
			return nil, err
		}

		return ChecksumResponse{
			id, base32.StdEncoding.EncodeToString(checksum),
		}, nil
	}
}
