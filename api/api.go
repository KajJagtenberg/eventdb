package api

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/constants"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/oklog/ulid"
	"github.com/tidwall/redcon"
)

func CommandHandler(s store.Store) Handler {
	return func(conn redcon.Conn, cmd redcon.Command) bool {
		query := strings.ToLower(string(cmd.Args[0]))

		switch query {
		default:
			conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")

			return false
		case "ping":
			conn.WriteString("PONG")
		case "quit":
			conn.WriteString("OK")
			conn.Close()
		case "version":
			conn.WriteString(constants.Version)
		case "size":
			size, err := s.Size()
			if err != nil {
				conn.WriteError(err.Error())
				return false
			}

			conn.WriteInt64(size)
		case "eventcount":
			count, err := s.EventCount()
			if err != nil {
				conn.WriteError(err.Error())
				return false
			}

			conn.WriteInt64(count)
		case "eventcountest":
			count, err := s.EventCountEstimate()
			if err != nil {
				conn.WriteError(err.Error())
				return false
			}

			conn.WriteInt64(count)
		case "streamcount":
			count, err := s.StreamCount()
			if err != nil {
				conn.WriteError(err.Error())
				return false
			}

			conn.WriteInt64(count)
		case "streamcountest":
			count, err := s.StreamCountEstimate()
			if err != nil {
				conn.WriteError(err.Error())
				return false
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
					return false
				}
			case 3:
				offset, err = ulid.Parse(string(cmd.Args[1]))
				if err != nil {
					conn.WriteError(err.Error())
					return false
				}

				limit, err = strconv.Atoi(string(cmd.Args[2]))
				if err != nil {
					conn.WriteError(err.Error())
					return false
				}
			default:
				conn.WriteError("Amount of arguments is not supported")
				return false
			}

			events, err := s.Log(offset, uint32(limit))
			if err != nil {
				conn.WriteError(err.Error())
				return false
			}

			result, err := json.Marshal(events)
			if err != nil {
				conn.WriteError(err.Error())
				return false
			}

			conn.WriteString(string(result))

		case "add":
			stream, err := uuid.ParseBytes(cmd.Args[1])
			if err != nil {
				conn.WriteError(err.Error())
				return false
			}

			version, err := strconv.ParseUint(string(cmd.Args[2]), 10, 32)
			if err != nil {
				conn.WriteError(err.Error())
				return false
			}

			var data []store.EventData

			if err := json.Unmarshal(cmd.Args[3], &data); err != nil {
				conn.WriteError(err.Error())
				return false
			}

			events, err := s.Add(stream, uint32(version), data)
			if err != nil {
				conn.WriteError(err.Error())
				return false
			}

			result, err := json.Marshal(events)
			if err != nil {
				conn.WriteError(err.Error())
				return false
			}

			conn.WriteString(string(result))

			// TODO: Add get
		}

		return true
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
