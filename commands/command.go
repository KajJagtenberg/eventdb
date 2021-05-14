package commands

import "errors"

var (
	ErrUnknownCommand = errors.New("unknown command")
)

type Command struct {
	Name string `json:"name"`
	Args []byte `json:"args"`
}

type CommandHandler = func(cmd Command) (interface{}, error)

type CommandRegistry struct {
	handlers map[string]CommandHandler
}

func (r *CommandRegistry) Register(name string, shorthand string, handler CommandHandler) {
	r.handlers[name] = handler
	r.handlers[shorthand] = handler
}

func (r *CommandRegistry) Handle(cmd Command) (interface{}, error) {
	handler := r.handlers[cmd.Name]
	if handler == nil {
		return nil, ErrUnknownCommand
	}

	return handler(cmd)
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		handlers: make(map[string]CommandHandler),
	}
}
