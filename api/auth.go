package api

import "github.com/tidwall/redcon"

func Authentication(password string, next redcon.HandlerFunc) redcon.HandlerFunc {
	return func(conn redcon.Conn, cmd redcon.Command) {
		next(conn, cmd)
	}
}
