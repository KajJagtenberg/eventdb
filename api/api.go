package api

import (
	"os"
	"strings"

	"github.com/kajjagtenberg/eventflowdb/constants"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/tidwall/redcon"
)

type ConnectionContext struct {
	authenticated bool
}

type Resp struct {
	store store.Store
}

func (r *Resp) CommandHandler(conn redcon.Conn, cmd redcon.Command) {
	ctx := assertContext(conn)

	query := strings.ToLower(string(cmd.Args[0]))

	if query == "auth" {
		if string(cmd.Args[1]) != os.Getenv("AUTH_KEY") {
			conn.WriteError("Unauthorized")
			return
		}

		ctx.authenticated = true

		conn.WriteString("Authenticated")
		return
	}

	if !ctx.authenticated {
		conn.WriteError("Unauthorized")
		return
	}

	switch query {
	default:
		conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
	case "ping":
		conn.WriteString("PONG")
	case "quit":
		conn.WriteString("OK")
		conn.Close()
	case "version":
		conn.WriteString(constants.Version)
	case "size":
		conn.WriteInt64(r.store.Size())

	case "eventcount":
		count, err := r.store.EventCount()
		if err != nil {
			conn.WriteError(err.Error())
			return
		}

		conn.WriteInt64(count)
	case "eventcountest":
		count, err := r.store.EventCountEstimate()
		if err != nil {
			conn.WriteError(err.Error())
			return
		}

		conn.WriteInt64(count)
	case "streamcount":
		count, err := r.store.StreamCount()
		if err != nil {
			conn.WriteError(err.Error())
			return
		}

		conn.WriteInt64(count)
	case "streamcountest":
		count, err := r.store.StreamCountEstimate()
		if err != nil {
			conn.WriteError(err.Error())
			return
		}

		conn.WriteInt64(count)
	}
}

func (r *Resp) AcceptHandler(conn redcon.Conn) bool {
	return true
}

func (r *Resp) ErrorHandler(conn redcon.Conn, err error) {

}

func NewResp(store store.Store) *Resp {
	return &Resp{store}
}

func assertContext(conn redcon.Conn) *ConnectionContext {
	ctx, ok := conn.Context().(*ConnectionContext)

	if !ok {
		ctx = &ConnectionContext{}

		conn.SetContext(ctx)
	}

	return ctx
}
