package commands_test

import (
	"testing"

	"github.com/kajjagtenberg/eventflowdb/commands"
	"github.com/kajjagtenberg/eventflowdb/constants"
	"github.com/kajjagtenberg/go-commando"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	handler := commands.VersionHandler()
	result, err := handler(commando.Command{
		Name: "version",
	})

	if err != nil {
		t.Fatal(err)
	}

	res, ok := result.(commands.VersionResponse)
	if !ok {
		t.Fatal("wrong cast")
	}

	assert := assert.New(t)
	assert.Equal(res.Version, constants.Version)
}
