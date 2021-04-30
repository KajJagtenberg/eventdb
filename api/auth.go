package api

import (
	"strings"

	"github.com/tidwall/redcon"
)

func Authentication(password string) Handler {
	return func(conn redcon.Conn, cmd redcon.Command) bool {
		if password == "" {
			return true
		}

		session := conn.Context().(*Session)

		if session.Authenticated {
			return true
		}

		query := strings.ToLower(string(cmd.Args[0]))

		if query != "auth" {
			return false
		}

		if string(cmd.Args[1]) != password {
			conn.WriteError("Unauthorized")
			return false
		}

		session.Authenticated = true

		conn.SetContext(session)

		conn.WriteString("OK")

		return false
	}
}
