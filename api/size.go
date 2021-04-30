package api

import "github.com/kajjagtenberg/eventflowdb/store"

func Size(store store.Store, c *Ctx) error {
	size, err := store.Size()
	if err != nil {
		return err
	}

	c.Conn.WriteInt64(size)

	return nil
}
