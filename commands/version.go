package commands

import "github.com/kajjagtenberg/eventflowdb/constants"

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
