package cluster

type FSM struct{}

func NewFSM() (*FSM, error) {
	return &FSM{}, nil
}
