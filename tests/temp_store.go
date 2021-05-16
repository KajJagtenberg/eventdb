package tests

import (
	"io/ioutil"
	"log"

	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
)

func CreateTempStore() (store.EventStore, error) {
	f, err := ioutil.TempFile("/tmp", "*")
	if err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}

	log.Println(f.Name())

	db, err := bbolt.Open(f.Name(), 0666, bbolt.DefaultOptions)
	if err != nil {
		return nil, err
	}

	return store.NewBoltStore(db, logrus.New())
}
