package main

import (
	"eventflowdb/compiler"
	"fmt"
	"log"
	"os"
	"strings"

	_ "embed"

	"github.com/chzyer/readline"
	"github.com/dop251/goja"
)

const (
	version = "0.0.1"
)

//go:embed shell.js
var shellSource string

var running = true

func main() {
	compiler, err := compiler.NewCompiler()
	if err != nil {
		log.Fatal(err)
	}

	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	vm.Set("version", func() {
		fmt.Println(version)
	})
	vm.Set("console", struct {
		Log interface{} `json:"log"`
	}{
		Log: func(v string) {
			fmt.Println(v)
		},
	})
	vm.RunString(shellSource)

	fmt.Println("EventflowDB Shell")

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     "/tmp/flowcli",
		EOFPrompt:       "exit",
		Stdout:          os.Stdout,
		Stderr:          os.Stderr,
		InterruptPrompt: "^C",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	for running {
		input, err := rl.Readline()
		if err != nil {
			panic(err)
		}

		tokens := strings.Split(input, " ")

		switch tokens[0] {
		case "exit":
			running = false
		case "version":
			fmt.Println(version)
		default:
			code, err := compiler.Compile(input)
			if err != nil {
				fmt.Println(err)
				continue
			}

			output, err := vm.RunString(code)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if !goja.IsUndefined(output) {
				fmt.Println(output)
			}
		}
	}

	rl.Clean()
	rl.Close()
}
