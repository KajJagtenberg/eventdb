package commands

type Command struct {
	Name string `json:"name"`
	Args []byte `json:"args"`
}

type CommandHandler interface {
	Handle(cmd Command) (interface{}, error)
}

type CommandRegistry struct {
	handlers map[string]CommandHandler
}

func (r *CommandRegistry) Register(name string, shorthand string, handler CommandHandler) {
	r.handlers[name] = handler
	r.handlers[shorthand] = handler
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		handlers: make(map[string]CommandHandler),
	}
}
