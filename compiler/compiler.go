package compiler

import (
	_ "embed"
	"errors"
	"log"

	"github.com/dop251/goja"
)

//go:embed js/babel.min.js
var babelSource string

//go:embed js/babel.env.min.js
var babelEnvSource string

//go:embed js/compiler.js
var compilerSource string

type Compiler struct {
	vm *goja.Runtime
}

func (compiler *Compiler) Compile(source string) (string, error) {
	v := compiler.vm.Get("compile")

	compile, ok := goja.AssertFunction(v)
	if !ok {
		return "", errors.New("Result is not a valid function")
	}

	v, err := compile(goja.Undefined(), compiler.vm.ToValue(source))
	if err != nil {
		return "", err
	}

	return source, nil
}

func NewCompiler() (*Compiler, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	vm.Set("console", struct {
		Log interface{} `json:"log"`
	}{
		Log: log.Println,
	})

	if _, err := vm.RunString(babelSource); err != nil {
		return nil, err
	}

	if _, err := vm.RunString(babelEnvSource); err != nil {
		return nil, err
	}

	if _, err := vm.RunString(compilerSource); err != nil {
		return nil, err
	}

	return &Compiler{vm}, nil
}
