package api

import "github.com/kajjagtenberg/eventflowdb/constants"

func Version(c *Ctx) error {
	c.Conn.WriteString(constants.Version)

	return nil
}
