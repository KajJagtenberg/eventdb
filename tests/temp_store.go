package tests

import (
	"io/ioutil"

	"go.etcd.io/bbolt"
)

func CreateTempDB() (*bbolt.DB, error) {
	f, err := ioutil.TempFile("/tmp", "*")
	if err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}

	return bbolt.Open(f.Name(), 0666, bbolt.DefaultOptions)
}
