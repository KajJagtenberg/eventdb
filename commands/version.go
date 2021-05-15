package commands

import (
	"github.com/kajjagtenberg/eventflowdb/constants"
	"github.com/kajjagtenberg/go-commando"
)

const (
	CMD_VERSION       = "version"
	CMD_VERSION_SHORT = "v"
)

type VersionResponse struct {
	Version string `json:"version"`
}

func VersionHandler() commando.CommandHandler {
	return func(cmd commando.Command) (interface{}, error) {
		return VersionResponse{
			Version: constants.Version,
		}, nil
	}
}
