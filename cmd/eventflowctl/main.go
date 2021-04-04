package main

import (
	"log"

	"github.com/chzyer/readline"
)

func main() {
	rl, err := readline.New("> ")
	if err != nil {
		log.Fatalf("Failed to open prompt: %v", err)
	}
	defer rl.Close()
}
