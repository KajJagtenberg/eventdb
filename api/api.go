package api

import (
	"errors"

	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/tidwall/redcon"
)

var (
	ErrUnknownCommand = errors.New("unknown command")
)

func CommandHandler(s store.Store) Handler {
	return func(c *Ctx) error {
		switch c.Command {

		case "streamcount":
			return StreamCount(s, c)
		case "streamcountest":
			return StreamCountEstimate(s, c)
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
