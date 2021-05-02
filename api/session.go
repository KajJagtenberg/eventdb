package api

type Session struct {
	Authenticated bool
}

func AssertSession() Handler {
	return func(c *Ctx) error {
		ctx := c.Conn.Context()

		if ctx == nil {
			ctx = &Session{
				Authenticated: false,
			}
		}

		c.Conn.SetContext(ctx)

		return c.Next()
	}
}
