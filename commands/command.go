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

type CommandDispatcher struct {
	handlers map[string]CommandHandler
}

func (r *CommandDispatcher) Register(name string, shorthand string, handler CommandHandler) {
	r.handlers[name] = handler
	r.handlers[shorthand] = handler
}

func (r *CommandDispatcher) Handle(cmd Command) (interface{}, error) {
	handler := r.handlers[cmd.Name]
	if handler == nil {
		return nil, ErrUnknownCommand
	}

	return handler(cmd)
}

func NewCommandDispatcher() *CommandDispatcher {
	return &CommandDispatcher{
		handlers: make(map[string]CommandHandler),
	}
}
