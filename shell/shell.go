package shell

import (
	"encoding/json"
	"eventflowdb/compiler"
	"eventflowdb/store"
	"fmt"
	"log"
	"os"

	_ "embed"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
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

	registry := new(require.Registry)
	registry.Enable(vm)

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
				Version:  int(record.Version),
				Type:     record.Type,
				Data:     record.Data,
				Metadata: record.Metadata,
				AddedAt:  record.AddedAt,
			}

			events = append(events, event)
		}

		return vm.ToValue(events), err
	})

	vm.Set("pwd", func() {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(dir)
	})

	// vm.Set("exports", goja.Object{})

	// vm.Set("require", func(src string) error {
	// 	log.Println(os.Getwd())

	// 	f, err := os.OpenFile("./shell/node_modules/"+src+".js", os.O_RDONLY, 0666)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer f.Close()

	// 	data, err := ioutil.ReadAll(f)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if _, err := vm.RunString(string(data)); err != nil {
	// 		return err
	// 	}

	// 	return nil
	// })

	vm.Set("project", func(projection map[string]func(state interface{}, event interface{})) {
		var offset ulid.ULID
		state := map[string]interface{}{}

		records, err := eventstore.LoadFromAll(offset, 10)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, record := range records {
			handler := projection[record.Type]

			if handler == nil {
				continue
			}

			var event interface{}

			if json.Unmarshal(record.Data, &event); err != nil {
				fmt.Println(err)
				return
			}

			handler(state, event)

			log.Println("Projection:", state)

			offset = record.ID
		}
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
