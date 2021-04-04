package shell

import (
	_ "embed"

	"github.com/dop251/goja"
)

var (
	//go:embed babel.min.js
	babelSource string

	//go:embed babel.preset.env.min.js
	babelEnvSource string
)

type Babel struct {
	vm *goja.Runtime
}

func (b *Babel) Compile(code string) (string, error) {
	if err := b.vm.Set("code", code); err != nil {
		return "", err
	}

	result, err := b.vm.RunString("compile()")
	return result.String(), err
}

func NewBabel() (*Babel, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	if _, err := vm.RunString(babelSource); err != nil {
		return nil, err
	}

	if _, err := vm.RunString(babelEnvSource); err != nil {
		return nil, err
	}

	if _, err := vm.RunString(`
		function compile() {
			return Babel.transform(code, {
				presets: ["env"]
			}).code
		}
	`); err != nil {
		return nil, err
	}

	return &Babel{vm}, nil
}
