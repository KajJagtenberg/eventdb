package api

import (
	"encoding/base64"
	"strings"

	"github.com/tidwall/redcon"
)

func Combine(handlers ...Handler) redcon.HandlerFunc {
	return func(conn redcon.Conn, cmd redcon.Command) {
		var args []byte

		if len(cmd.Args) > 2 {
			conn.WriteError("Amount of arguments cannot be more than 2")
			return
		}

		if len(cmd.Args) == 2 {
			var err error
			args, err = base64.StdEncoding.DecodeString(string(cmd.Args[1]))
			if err != nil {
				conn.WriteError(err.Error())
				return
			}
		}

		ctx := &Ctx{
			Conn:    conn,
			Command: strings.ToLower(string(cmd.Args[0])),
			Args:    args,
		}
		ctx.next = handlers[1:]

		err := handlers[0](ctx)
		if err != nil {
			conn.WriteError(err.Error())
		}
	}
}
