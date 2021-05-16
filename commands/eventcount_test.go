package commands_test

import (
	"log"
	"os"
	"testing"

	"github.com/kajjagtenberg/eventflowdb/commands"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/kajjagtenberg/eventflowdb/tests"
	"github.com/kajjagtenberg/go-commando"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestEventCount(t *testing.T) {
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

	handler := commands.EventCountHandler(store)

	cmd := commando.Command{}

	result, err := handler(cmd)
	if err != nil {
		t.Fatal(err)
	}

	res, ok := result.(commands.EventCountResponse)
	if !ok {
		t.Fatal("Wrong cast")
	}

	assert := assert.New(t)
	assert.Equal(res.Count, int64(0))
}
