package api

import "github.com/tidwall/redcon"

type Handler = func(conn redcon.Conn, cmd redcon.Command) bool

func Combine(handlers ...Handler) redcon.HandlerFunc {
	return func(conn redcon.Conn, cmd redcon.Command) {
		for _, handler := range handlers {
			if !handler(conn, cmd) {
				break
			}
		}
	}
}
