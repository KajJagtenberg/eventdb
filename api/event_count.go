package api

import "github.com/KajJagtenberg/eventflowdb/store"

func EventCount(store store.Store, c *Ctx) error {
	count, err := store.EventCount()
	if err != nil {
		return err
	}

	c.Conn.WriteInt64(count)

	return nil
}

func EventCountEstimate(store store.Store, c *Ctx) error {
	count, err := store.EventCountEstimate()
	if err != nil {
		return err
	}

	c.Conn.WriteInt64(count)

	return nil
}
