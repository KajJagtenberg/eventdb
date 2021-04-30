package api

import (
	"errors"

	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/tidwall/redcon"
)

func CommandHandler(s store.Store) Handler {
	return func(c *Ctx) error {
		return errors.New("Unknown command")
	}
}

func AcceptHandler() func(conn redcon.Conn) bool {
	return func(conn redcon.Conn) bool {
		return true
	}
}

func ErrorHandler() func(conn redcon.Conn, err error) {
	return func(conn redcon.Conn, err error) {

	}
}
