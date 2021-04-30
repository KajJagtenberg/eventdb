package api

import (
	"github.com/tidwall/redcon"
)

func Authentication(password string) Handler {
	return func(conn redcon.Conn, cmd redcon.Command) bool {
		conn.WriteError("Unauthorized")

		return false
	}
}
