package api

import (
	"errors"
)

var (
	ErrUnauthorized         = errors.New("Unauthorized")
	ErrAlreadyAuthenticated = errors.New("Already authenticated")
)

func Authentication(password string) Handler {
	return func(c *Ctx) error {
		session := c.Conn.Context().(*Session)

		if session.Authenticated {
			if c.Command == "auth" {
				return ErrAlreadyAuthenticated
			} else {
				return c.Next()
			}
		} else {
			if password == "" {
				return c.Next()
			}

			if c.Command != "auth" {
				return ErrUnauthorized
			}

			if string(c.Args[0]) != password {
				return ErrUnauthorized
			}

			session.Authenticated = true

			c.Conn.SetContext(session)

			c.Conn.WriteString("OK")

			return nil
		}
	}
}
