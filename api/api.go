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
	return func(ctx *Ctx) error {
		switch ctx.Command {
		case "ping":
			return Ping(ctx)
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
