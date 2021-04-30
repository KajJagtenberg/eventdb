package api

import (
	"github.com/tidwall/redcon"
)

type Handler = func(ctx *Ctx) error

type Ctx struct {
	Conn    redcon.Conn
	Command string
	Args    [][]byte

	next interface{}
}

func (c *Ctx) Next() error {
	next := c.next.([]Handler)
	if len(next) == 0 {
		return nil
	}

	h := next[0]
	c.next = next[1:]
	return h(c)
}
