package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/kajjagtenberg/eventflowdb/api"
	"google.golang.org/grpc"
)

const (
	Version = "0.1.0"
)

func main() {
	rl, err := readline.New("> ")
	if err != nil {
		log.Fatalf("Failed to open prompt: %v", err)
	}
	defer rl.Close()

	conn, err := grpc.Dial("127.0.0.1:6543", grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	defer conn.Close()

	client := api.NewShellServiceClient(conn)
	exec, err := client.Execute(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

Loop:
	for {
		line, err := rl.Readline()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read line: %v", err)
		}

		parts := strings.Split(line, " ")

		switch strings.ToLower(parts[0]) {
		case "exit":
			break Loop
		case "version":
			fmt.Println(Version)
		default:
			if err := exec.Send(&api.ShellRequest{
				Body: line,
			}); err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}

			response, err := exec.Recv()
			if err != nil {
				fmt.Println(err)
			}

			if body := response.Body; body != "" {
				fmt.Println(body)
			}
		}
	}

	fmt.Println("Goodbye!")
}
