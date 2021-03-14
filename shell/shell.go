package shell

import (
	"eventflowdb/compiler"
	"log"
	"os/exec"

	_ "embed"

	"github.com/dop251/goja"
)

//go:embed shell.js
var source string

type Shell struct {
	vm *goja.Runtime
}

func NewShell() (*Shell, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	vm.Set("exec", func(name string, arg ...string) (string, error) {
		cmd := exec.Command(name, arg...)
		output, err := cmd.Output()
		return string(output), err
	})
	vm.Set("console", struct {
		Log interface{} `json:"log"`
	}{
		Log: log.Println,
	})

	compiled, err := compiler.Compile(source)
	if err != nil {
		return nil, err
	}

	if _, err := vm.RunString(compiled); err != nil {
		return nil, err
	}

	return &Shell{vm}, nil
}
