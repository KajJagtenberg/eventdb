package resp

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/kajjagtenberg/eventflowdb/commands"
	"github.com/kajjagtenberg/go-commando"
	"github.com/tidwall/redcon"
)

type Session struct {
	Authenticated bool
}

func CommandHandler(dispatcher *commando.CommandDispatcher, password string) func(conn redcon.Conn, cmd redcon.Command) {
	return func(conn redcon.Conn, cmd redcon.Command) {
		if len(cmd.Args) == 0 {
			conn.WriteError("no command specified")
			return
		}

		var session *Session

		if value := conn.Context(); value != nil {
			session = value.(*Session)
		} else {
			session = &Session{
				Authenticated: false,
			}
		}

		if !session.Authenticated && len(password) > 0 {
			if strings.ToLower(string(cmd.Args[0])) != "auth" {
				conn.WriteError("unauthorized")
				return
			}

			if len(cmd.Args) == 1 {
				conn.WriteError("unauthorized")
				return
			}

			if string(cmd.Args[1]) != password {
				conn.WriteError("unauthorized")
				return
			}

			session.Authenticated = true

			conn.SetContext(session)

			conn.WriteString("OK")
			return
		}

		var args []byte

		if len(cmd.Args) > 1 {
			args = cmd.Args[1]
		}

		result, err := dispatcher.Handle(commando.Command{
			Name: strings.ToLower(string(cmd.Args[0])),
			Args: args,
		})
		if err != nil {
			conn.WriteError(err.Error())
			return
		}

		switch r := result.(type) {
		case commands.UptimeResponse:
			conn.WriteArray(2)
			conn.WriteInt64(r.Uptime)
			conn.WriteString(r.Human)

		case commands.VersionResponse:
			conn.WriteString(r.Version)

		case commands.AddResponse:
			conn.WriteArray(len(r.Events))

			for _, event := range r.Events {
				v, err := json.Marshal(&event)
				if err != nil {
					conn.WriteError(err.Error())
					return
				}

				conn.WriteString(string(v))
			}

		case commands.ChecksumResponse:
			conn.WriteArray(2)
			conn.WriteString(r.ID.String())
			conn.WriteString(r.Checksum)

		case commands.EventCountResponse:
			conn.WriteInt64(r.Count)

		case commands.GetResponse:
			conn.WriteArray(len(r.Events))

			for _, event := range r.Events {
				v, err := json.Marshal(&event)
				if err != nil {
					conn.WriteError(err.Error())
					return
				}

				conn.WriteString(string(v))
			}

		case commands.GetAllResponse:
			conn.WriteArray(len(r.Events))

			for _, event := range r.Events {
				v, err := json.Marshal(&event)
				if err != nil {
					conn.WriteError(err.Error())
					return
				}

				conn.WriteString(string(v))
			}

		case commands.PingResponse:
			conn.WriteString(r.Message)

		case commands.SizeResponse:
			conn.WriteArray(2)
			conn.WriteInt64(r.Size)
			conn.WriteString(r.Human)

		case commands.StreamCountResponse:
			conn.WriteInt64(r.Count)

		case commands.ListStreamsResponse:
			conn.WriteArray(len(r.Streams))

			for _, stream := range r.Streams {
				value, err := json.Marshal(stream)
				if err != nil {
					conn.WriteError(err.Error())
					return
				}

				conn.WriteString(string(value))
			}
		default:
			log.Println("No known result")
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
