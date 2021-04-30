package api

import (
	"github.com/tidwall/redcon"
)

func Combine(handlers ...Handler) redcon.HandlerFunc {
	return func(conn redcon.Conn, cmd redcon.Command) {

	}
}
