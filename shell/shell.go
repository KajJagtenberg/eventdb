package shell

import (
	"log"
	"time"

	_ "embed"

	"github.com/dop251/goja"
	"github.com/hashicorp/raft"
	"github.com/kajjagtenberg/eventflowdb/constants"
	"github.com/kajjagtenberg/eventflowdb/persistence"
	"github.com/oklog/ulid"
)

//go:embed shell.js
var runtime string

type Shell struct {
	vm    *goja.Runtime
	babel *Babel
}

func (shell *Shell) Execute(src string) (string, error) {
	compiled, err := shell.babel.Compile(src)
	if err != nil {
		return "", nil
	}

	value, err := shell.vm.RunString(compiled)
	if err != nil {
		return "", err
	}

	body := value.String()

	if body == "use strict" {
		body = ""
	}

	return body, nil
}

func NewShell(raftServer *raft.Raft, persistence *persistence.Persistence) *Shell {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	vm.Set("version", func() string {
		return constants.Version
	})

	vm.Set("clusterLeader", func() string {
		return string(raftServer.Leader())
	})

	vm.Set("clusterStats", func() interface{} {
		return raftServer.Stats()
	})

	vm.Set("logEvents", func() (interface{}, error) {
		events, err := persistence.Log(ulid.ULID{}, 100)
		return events, err
	})

	vm.Set("joinCluster", func(id string, address string) error {
		future := raftServer.AddVoter(raft.ServerID(id), raft.ServerAddress(address), 0, time.Second*5)
		return future.Error()
	})

	vm.Set("clusterSize", func() int {
		return len(raftServer.GetConfiguration().Configuration().Servers)
	})

	babel := NewBabel()

	shell := &Shell{vm, babel}

	if _, err := shell.Execute(runtime); err != nil {
		log.Fatalf("Cannot execute shell runtime. This is a bug: %v", err)
	}

	return shell
}
