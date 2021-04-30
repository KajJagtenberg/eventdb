package api

import (
	"github.com/tidwall/redcon"
)

type Handler = func(ctx *Ctx) error

type Ctx struct {
	Conn    redcon.Conn
	Command string
	Args    [][]byte

	next []interface{}
}

func (c *Ctx) Next() error {
	h := c.next[0].(Handler)
	c.next = c.next[1:]

	return h(c)
}
