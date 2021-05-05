package api

import "github.com/KajJagtenberg/eventflowdb/constants"

func Version(c *Ctx) error {
	c.Conn.WriteString(constants.Version)

	return nil
}
