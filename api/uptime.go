package api

import "time"

var (
	start = time.Now()
)

func Uptime(c *Ctx) error {
	uptime := time.Since(start)

	c.Conn.WriteString(uptime.String())

	return nil
}
