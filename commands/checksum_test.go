package commands_test

import (
	"encoding/base32"
	"testing"

	"github.com/kajjagtenberg/eventflowdb/commands"
	"github.com/kajjagtenberg/eventflowdb/tests"
	"github.com/kajjagtenberg/go-commando"
	"github.com/oklog/ulid"
	"github.com/stretchr/testify/assert"
)

func TestChecksum(t *testing.T) {
	store, err := tests.CreateTempStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	handler := commands.ChecksumHandler(store)

	cmd := commando.Command{}

	result, err := handler(cmd)
	if err != nil {
		t.Fatal(err)
	}

	res, ok := result.(commands.ChecksumResponse)
	if !ok {
		t.Fatal("Wrong cast")
	}

	assert := assert.New(t)
	assert.Equal(res.ID, ulid.ULID{})
	assert.Equal(res.Checksum, "AAAAAAA=")

	_, err = base32.StdEncoding.DecodeString(res.Checksum)
	assert.Equal(err, nil)
}
