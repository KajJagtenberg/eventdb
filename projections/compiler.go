package projections

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/dop251/goja"
)

type Compiler struct {
	vm *goja.Runtime
}

func NewCompiler() (*Compiler, error) {
	vm := goja.New()

	if err := loadBabel(vm); err != nil {
		return nil, err
	}

	return &Compiler{vm}, nil
}

func (c *Compiler) Compile(code string) (string, error) {
	c.vm.Set("input", code)

	if _, err := c.vm.RunString(`var output = Babel.transform(input, {presets: ["es2015"]}).code;`); err != nil {
		return "", err
	}

	return c.vm.Get("output").String(), nil
}

func loadBabel(vm *goja.Runtime) error {
	file, err := os.OpenFile("projections/babel.min.js", os.O_RDONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	src, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	_, err = vm.RunString(string(src))
	return err
}
