package commands

const (
	CMD_PING       = "ping"
	CMD_PING_SHORT = "p"
)

type PingResponse struct {
	Message string `json:"message"`
}

func PingHandler() CommandHandler {
	return func(cmd Command) (interface{}, error) {
		return PingResponse{
			Message: "PONG",
		}, nil
	}
}
