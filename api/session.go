package api

import (
	"github.com/tidwall/redcon"
)

type Session struct {
	Authenticated bool
}

func AssertSession() Handler {
	return func(conn redcon.Conn, cmd redcon.Command) bool {
		ctx := conn.Context()

		if ctx == nil {
			ctx = &Session{
				Authenticated: false,
			}
		}

		conn.SetContext(ctx)

		return true
	}
}
