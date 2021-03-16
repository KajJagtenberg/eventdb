package shell

import (
	"eventflowdb/compiler"
	"eventflowdb/store"
	"fmt"
	"log"

	_ "embed"

	"github.com/dop251/goja"
	"github.com/oklog/ulid"
)

//go:embed shell.js
var source string

type Shell struct {
	vm *goja.Runtime
}

func (shell *Shell) Run(code string) (string, error) {
	code, err := compiler.Compile(code)
	if err != nil {
		return "", err
	}

	output, err := shell.vm.RunString(code)
	if err != nil {
		return "", err
	}

	if !goja.IsUndefined(output) && !output.Equals(shell.vm.ToValue("use strict")) {
		return output.String(), err
	}

	return "", nil
}

func NewShell(eventstore *store.EventStore) (*Shell, error) {

	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	vm.Set("console", struct {
		Log interface{} `json:"log"`
	}{
		Log: func(v string) {
			fmt.Println(v)
		},
	})
	vm.Set("loadFromAll", func() (goja.Value, error) {
		records, err := eventstore.LoadFromAll(ulid.ULID{}, 10)

		var events []Event

		for _, record := range records {
			event := Event{
				ID:       record.ID.String(),
				Stream:   record.Stream.String(),
				Version:  record.Version,
				Type:     record.Type,
				Data:     record.Data,
				Metadata: record.Metadata,
				AddedAt:  record.AddedAt,
			}

			events = append(events, event)
		}

		return vm.ToValue(events), err
	})

	freeze, ok := goja.AssertFunction(vm.Get("Object").ToObject(vm).Get("freeze"))
	if !ok {
		log.Fatal("Result is not a valid function")
	}

	for _, key := range vm.GlobalObject().Keys() {
		freeze(goja.Undefined(), vm.ToValue(key))
	}

	compiled, err := compiler.Compile(source)
	if err != nil {
		return nil, err
	}

	if _, err := vm.RunString(compiled); err != nil {
		return nil, err
	}

	return &Shell{vm}, nil
}
