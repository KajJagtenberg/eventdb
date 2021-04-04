package shell

import (
	_ "embed"

	"github.com/dop251/goja"
)

var (
	//go:embed runtime.js
	runtime string
)

type Shell struct {
	vm *goja.Runtime
}

func NewShell() (*Shell, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	if _, err := vm.RunString(runtime); err != nil {
		return nil, err
	}

	return &Shell{vm}, nil
}
