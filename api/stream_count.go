package api

import "github.com/KajJagtenberg/eventflowdb/store"

func StreamCount(store store.Store, c *Ctx) error {
	count, err := store.StreamCount()
	if err != nil {
		return err
	}

	c.Conn.WriteInt64(count)

	return nil
}

func StreamCountEstimate(store store.Store, c *Ctx) error {
	count, err := store.StreamCountEstimate()
	if err != nil {
		return err
	}

	c.Conn.WriteInt64(count)

	return nil
}
