package api

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/kajjagtenberg/eventflowdb/constants"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/oklog/ulid"
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
		size, err := r.store.Size()
		if err != nil {
			conn.WriteError(err.Error())
			return
		}

		conn.WriteInt64(size)
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
	case "log":
		var offset ulid.ULID
		var limit int
		var err error

		switch len(cmd.Args) {
		case 1:
			offset = ulid.ULID{}
			limit = 10
		case 2:
			offset, err = ulid.Parse(string(cmd.Args[1]))
			if err != nil {
				conn.WriteError(err.Error())
				return
			}
		case 3:
			offset, err = ulid.Parse(string(cmd.Args[1]))
			if err != nil {
				conn.WriteError(err.Error())
				return
			}

			limit, err = strconv.Atoi(string(cmd.Args[2]))
			if err != nil {
				conn.WriteError(err.Error())
				return
			}
		default:
			conn.WriteError("Amount of arguments is not supported")
			return
		}

		events, err := r.store.Log(offset, uint32(limit))
		if err != nil {
			conn.WriteError(err.Error())
			return
		}

		conn.WriteArray(len(events))

		for _, event := range events {
			raw, err := json.Marshal(event)
			if err != nil {
				conn.WriteError(err.Error())
				return
			}

			conn.WriteString(string(raw))
		}
	}
}

func (r *Resp) AcceptHandler(conn redcon.Conn) bool {
	return true
}

func (r *Resp) ErrorHandler(conn redcon.Conn, err error) {}

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
