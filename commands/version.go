package commands

import "github.com/kajjagtenberg/eventflowdb/constants"

const (
	CMD_VERSION       = "version"
	CMD_VERSION_SHORT = "v"
)

type VersionResponse struct {
	Version string `json:"version"`
}

func VersionHandler() CommandHandler {
	return func(cmd Command) (interface{}, error) {
		return VersionResponse{
			Version: constants.Version,
		}, nil
	}
}
