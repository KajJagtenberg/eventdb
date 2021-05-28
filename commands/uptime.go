package commands

import (
	"time"

	"github.com/kajjagtenberg/go-commando"
)

const (
	CMD_UPTIME       = "uptime"
	CMD_UPTIME_SHORT = "up"
)

type UptimeResponse struct {
	Uptime int64  `json:"uptime"`
	Human  string `json:"uptime_human"`
}

func UptimeHandler() commando.CommandHandler {
	start := time.Now()

	return func(cmd commando.Command) (interface{}, error) {
		uptime := time.Since(start)
		return UptimeResponse{
			Uptime: uptime.Milliseconds(),
			Human:  uptime.String(),
		}, nil
	}
}

func SetupUptimeHandler(dispatcher *commando.CommandDispatcher) {
	dispatcher.Register(CMD_UPTIME, CMD_UPTIME_SHORT, UptimeHandler())
}
