package commands

import "time"

type UptimeResponse struct {
	Uptime time.Duration `json:"uptime"`
	Human  string        `json:"uptime_human"`
}

func UptimeHandler() CommandHandler {
	start := time.Now()

	return func(cmd Command) (interface{}, error) {
		uptime := time.Since(start)
		return UptimeResponse{
			Uptime: uptime,
			Human:  uptime.String(),
		}, nil
	}
}
