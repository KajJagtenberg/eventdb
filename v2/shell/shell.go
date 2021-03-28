package shell

import (
	"log"

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

func NewShell(raft *raft.Raft, persistence *persistence.Persistence) *Shell {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	vm.Set("version", func() string {
		return constants.Version
	})

	vm.Set("leader", func() string {
		return string(raft.Leader())
	})

	vm.Set("stats", func() interface{} {
		return raft.Stats()
	})

	vm.Set("log", func() (interface{}, error) {
		events, err := persistence.Log(ulid.ULID{}, 100)

		log.Println(events[0].Type)

		return events, err
	})

	babel := NewBabel()

	shell := &Shell{vm, babel}

	if _, err := shell.Execute(runtime); err != nil {
		log.Fatalf("Cannot execute shell runtime. This is a bug: %v", err)
	}

	return shell
}
