package api

import (
	"errors"

	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/tidwall/redcon"
)

var (
	ErrUnknownCommand = errors.New("Unknown command")
)

func CommandHandler(s store.Store) Handler {
	return func(c *Ctx) error {
		switch c.Command {
		case "ping":
			return Ping(c)
		case "quit":
			return Quit(c)
		case "size":
			return Size(s, c)
		case "eventcount":
			return EventCount(s, c)
		default:
			return ErrUnknownCommand
		}
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
