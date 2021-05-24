package commands_test

import (
	"testing"

	"github.com/kajjagtenberg/eventflowdb/commands"
	"github.com/kajjagtenberg/go-commando"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	handler := commands.PingHandler()
	result, err := handler(commando.Command{
		Name: "ping",
	})
	if err != nil {
		t.Fatal(err)
	}

	res, ok := result.(commands.PingResponse)
	if !ok {
		t.Fatal("wrong cast")
	}

	assert := assert.New(t)
	assert.Equal(res.Message, "PONG")
}
