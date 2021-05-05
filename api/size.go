package api

import (
	"github.com/KajJagtenberg/eventflowdb/si"
	"github.com/KajJagtenberg/eventflowdb/store"
)

func Size(store store.Store, c *Ctx) error {
	size, err := store.Size()
	if err != nil {
		return err
	}

	human := si.ByteCountSI(size)

	c.Conn.WriteArray(2)
	c.Conn.WriteInt64(size)
	c.Conn.WriteString(human)

	return nil
}
