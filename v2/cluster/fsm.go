package cluster

import (
	"io"
	"log"
	"strings"
	"sync"

	"github.com/hashicorp/raft"
)

type FSM struct {
	state map[string]string
	lock  sync.RWMutex
}

func (fsm *FSM) Apply(applyLog *raft.Log) interface{} {
	switch applyLog.Type {
	case raft.LogCommand:

		cmd := string(applyLog.Data)

		tokens := strings.Split(cmd, " ")

		switch tokens[0] {
		case "SET":
			subtokens := strings.Split(tokens[1], "=")
			key := subtokens[0]
			value := subtokens[1]

			fsm.Set(key, value)
		default:
			log.Println("Unknown command")
		}

	default:
		log.Println("Type:", applyLog.Type)
	}

	return nil
}

func (fsm *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (fsm *FSM) Set(key string, value string) {
	fsm.lock.Lock()
	defer fsm.lock.Unlock()

	fsm.state[key] = value
}

func (fsm *FSM) Get(key string) string {
	fsm.lock.Lock()
	defer fsm.lock.Unlock()

	return fsm.state[key]
}

func (fsm *FSM) GetState() map[string]string {
	fsm.lock.Lock()
	defer fsm.lock.Unlock()

	return fsm.state
}

func (fsm *FSM) Restore(io.ReadCloser) error {
	return nil
}

func NewFSM() (*FSM, error) {
	return &FSM{state: make(map[string]string)}, nil
}
