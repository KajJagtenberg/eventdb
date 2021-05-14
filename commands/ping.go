package commands

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
