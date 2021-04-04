package shell

import (
	_ "embed"

	"github.com/dop251/goja"
	"github.com/kajjagtenberg/eventflowdb/constants"
)

var (
	//go:embed runtime.js
	runtime string
	babel   *Babel
)

type Shell struct {
	vm *goja.Runtime
}

func (shell *Shell) Execute(code string) (string, error) {
	if babel == nil {
		var err error
		babel, err = NewBabel()
		if err != nil {
			return "", err
		}
	}

	compiled, err := babel.Compile(code)
	if err != nil {
		if err != nil {
			return "", err
		}
	}

	value, err := shell.vm.RunString(compiled)
	if err != nil {
		return err.Error(), nil
	}

	result := value.String()

	if result == "use strict" {
		return "", nil
	}

	return result, nil
}

func NewShell() (*Shell, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	if _, err := vm.RunString(runtime); err != nil {
		return nil, err
	}

	vm.Set("version", func() string {
		return constants.Version
	})

	return &Shell{vm}, nil
}
