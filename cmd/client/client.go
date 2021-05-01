package main

import (
	"log"

	v1 "github.com/kajjagtenberg/eventflowdb/client"
)

func main() {
	client := v1.NewClient(&v1.Config{
		Address: "127.0.0.1:6543",
	})

	log.Println(client.Size())
}
