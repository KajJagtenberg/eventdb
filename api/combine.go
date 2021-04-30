package api

import (
	"strings"

	"github.com/tidwall/redcon"
)

func Combine(handlers ...Handler) redcon.HandlerFunc {
	return func(conn redcon.Conn, cmd redcon.Command) {
		ctx := &Ctx{
			Conn:    conn,
			Command: strings.ToLower(string(cmd.Args[0])),
			Args:    cmd.Args[1:],
		}
		ctx.next = handlers

		err := handlers[0](ctx)
		if err != nil {
			conn.WriteError(err.Error())
		}
	}
}
