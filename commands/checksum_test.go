package commands_test

import (
	"encoding/base32"
	"log"
	"os"
	"testing"

	"github.com/kajjagtenberg/eventflowdb/commands"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/kajjagtenberg/eventflowdb/tests"
	"github.com/kajjagtenberg/go-commando"
	"github.com/oklog/ulid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestChecksum(t *testing.T) {
	db, err := tests.CreateTempDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove(db.Path())

	store, err := store.NewBoltStore(db, logrus.New())
	if err != nil {
		log.Fatal(err)
	}

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
