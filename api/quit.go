package api

func Quit(c *Ctx) error {
	return c.Conn.Close()
}
