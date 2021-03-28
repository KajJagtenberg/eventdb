package shell

import (
	"log"

	_ "embed"

	"github.com/dop251/goja"
)

//go:embed babel.min.js
var source string

//go:embed babel.env.min.js
var preset string

type Babel struct {
	vm *goja.Runtime
}

func (babel *Babel) Compile(source string) (string, error) {
	if err := babel.vm.Set("source", source); err != nil {
		return "", err
	}

	value, err := babel.vm.RunString(`Babel.transform(source,{presets:["env"]}).code;`)
	if err != nil {
		return "", nil
	}

	return value.String(), nil
}

func NewBabel() *Babel {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	if _, err := vm.RunString(source); err != nil {
		log.Fatalf("Failed to run load babel: %v", err)
	}

	if _, err := vm.RunString(preset); err != nil {
		log.Fatalf("Failed to run load babel preset: %v", err)
	}

	return &Babel{vm}
}
