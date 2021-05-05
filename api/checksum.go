package api

import (
	"encoding/base32"

	"github.com/KajJagtenberg/eventflowdb/store"
)

func Checksum(s store.Store, c *Ctx) error {
	id, checksum, err := s.Checksum()
	if err != nil {
		return err
	}

	result := base32.StdEncoding.EncodeToString(checksum)

	c.Conn.WriteArray(2)
	c.Conn.WriteString(id.String())
	c.Conn.WriteString(result)

	return nil
}
