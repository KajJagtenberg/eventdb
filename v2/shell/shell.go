package shell

import (
	"log"

	_ "embed"

	"github.com/dop251/goja"
	"github.com/hashicorp/raft"
	babel "github.com/jvatic/goja-babel"
	"github.com/kajjagtenberg/eventflowdb/constants"
)

//go:embed shell.js
var runtime string

type Shell struct {
	vm *goja.Runtime
}

func (shell *Shell) Execute(src string) (string, error) {
	compiled, err := babel.TransformString(src, map[string]interface{}{
		"presets": []string{
			"env",
		},
	})
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

func NewShell(raft *raft.Raft) *Shell {
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

	shell := &Shell{vm}

	if _, err := shell.Execute(runtime); err != nil {
		log.Fatalf("Cannot execute shell runtime. This is a bug: %v", err)
	}

	return shell
}
