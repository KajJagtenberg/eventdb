package projections

import (
	_ "embed"

	"github.com/dop251/goja"
)

//go:embed js/babel.min.js
var babelSource string

type Compiler struct {
	vm *goja.Runtime
}

func (compiler *Compiler) Compile(source string) (string, error) {
	if err := compiler.vm.Set("src", source); err != nil {
		return "", err
	}

	value, err := compiler.vm.RunString(`
		Babel.transform(src, {}).code;
	`)
	if err != nil {
		return "", err
	}

	result := value.Export().(string)

	return result, nil
}

func NewCompiler() (*Compiler, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	if _, err := vm.RunString(babelSource); err != nil {
		return nil, err
	}

	return &Compiler{vm}, nil
}
