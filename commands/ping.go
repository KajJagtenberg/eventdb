package commands

import "github.com/kajjagtenberg/go-commando"

const (
	CMD_PING       = "ping"
	CMD_PING_SHORT = "p"
)

type PingResponse struct {
	Message string `json:"message"`
}

func PingHandler() commando.CommandHandler {
	return func(cmd commando.Command) (interface{}, error) {
		return PingResponse{
			Message: "PONG",
		}, nil
	}
}

func SetupPingHandler(dispatcher *commando.CommandDispatcher) {
	dispatcher.Register(CMD_PING, CMD_PING_SHORT, PingHandler())
}
