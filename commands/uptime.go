package commands

import (
	"time"

	"github.com/kajjagtenberg/go-commando"
)

const (
	CMD_UPTIME       = "uptime"
	CMD_UPTIME_SHORT = "u"
)

type UptimeResponse struct {
	Uptime time.Duration `json:"uptime"`
	Human  string        `json:"uptime_human"`
}

func UptimeHandler() commando.CommandHandler {
	start := time.Now()

	return func(cmd commando.Command) (interface{}, error) {
		uptime := time.Since(start)
		return UptimeResponse{
			Uptime: uptime,
			Human:  uptime.String(),
		}, nil
	}
}
