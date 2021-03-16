package main

import (
	"eventflowdb/compiler"
	"eventflowdb/shell"
	"eventflowdb/store"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	_ "embed"

	"github.com/chzyer/readline"
	"go.etcd.io/bbolt"
)

const (
	version = "0.0.1"
)

var running = true

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, err := bbolt.Open("events.db", 0600, nil)
	check(err)
	defer db.Close()

	eventstore, err := store.NewEventStore(db)
	check(err)

	shell, err := shell.NewShell(eventstore)
	check(err)

	fmt.Println("EventflowDB Shell")

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     "/tmp/flowcli",
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
		if err == readline.ErrInterrupt {
			break
		}
		if err == io.EOF {
			break
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

			output, err := shell.Run(code)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if len(output) > 0 {
				fmt.Println(output)
			}
		}
	}

	rl.Clean()
	rl.Close()
}
