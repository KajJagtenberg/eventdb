package shell

import (
	_ "embed"
	"log"

	"github.com/dop251/goja"
	"github.com/kajjagtenberg/eventflowdb/constants"
	"github.com/kajjagtenberg/eventflowdb/store"
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
	compiled, err := babel.Compile(code)
	if err != nil {
		return err.Error(), nil
	}

	value, err := shell.vm.RunString(compiled)
	if err != nil {
		return err.Error(), nil
	}

	return value.String(), nil
}

func NewShell(store store.Store) (*Shell, error) {
	if babel == nil {
		var err error
		babel, err = NewBabel()
		if err != nil {
			return nil, err
		}
	}

	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	compiledRuntime, err := babel.Compile(runtime)
	if err != nil {
		log.Println("Runtime compilation failed:", err)
		return nil, err
	}

	if _, err := vm.RunString(compiledRuntime); err != nil {
		return nil, err
	}

	if err := vm.Set("version", func() string {
		return constants.Version
	}); err != nil {
		return nil, err
	}

	database := struct {
		StreamCount interface{} `json:"streamCount"`
		EventCount  interface{} `json:"eventCount"`
		Size        interface{} `json:"size"`
		Log         interface{} `json:"log"`
	}{
		StreamCount: store.StreamCount,
		EventCount:  store.EventCount,
		Size:        store.Size,
		Log:         store.Log,
	}

	if err := vm.Set("db", database); err != nil {
		return nil, err
	}

	return &Shell{vm}, nil
}
