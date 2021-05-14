package resp

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/kajjagtenberg/eventflowdb/commands"
	"github.com/tidwall/redcon"
)

func CommandHandler(dispatcher *commands.CommandDispatcher) func(conn redcon.Conn, cmd redcon.Command) {
	return func(conn redcon.Conn, cmd redcon.Command) {
		var args []byte

		if len(cmd.Args) > 1 {
			args = cmd.Args[1]
		}

		result, err := dispatcher.Handle(commands.Command{
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
			conn.WriteInt64(r.Uptime.Milliseconds())
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
