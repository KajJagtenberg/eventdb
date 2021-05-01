package api

import (
	"encoding/base32"

	"github.com/kajjagtenberg/eventflowdb/store"
)

func Checksum(s store.Store, c *Ctx) error {
	checksum, err := s.Checksum()
	if err != nil {
		return err
	}

	result := base32.StdEncoding.EncodeToString(checksum)

	c.Conn.WriteString(result)

	return nil
}
