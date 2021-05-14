package commands

type Command struct {
	Name string `json:"name"`
	Args []byte `json:"args"`
}

type CommandHandler interface {
	Handle(cmd Command) (interface{}, error)
}
