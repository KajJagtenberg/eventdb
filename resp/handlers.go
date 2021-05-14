package resp

import (
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
