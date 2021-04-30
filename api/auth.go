package api

import (
	"strings"

	"github.com/tidwall/redcon"
)

func Authentication(password string) Handler {
	return func(conn redcon.Conn, cmd redcon.Command) bool {
		session := conn.Context().(*Session)

		query := strings.ToLower(string(cmd.Args[0]))

		if session.Authenticated {
			if query == "auth" {
				conn.WriteError("Already authenticated")
				return false
			} else {
				return true
			}
		} else {
			if password == "" {
				return true
			}

			if query != "auth" {
				conn.WriteError("Unauthorized")
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
}
