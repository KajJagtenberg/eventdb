package fsm

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/hashicorp/raft"
)

type badgerFSM struct {
	db *badger.DB
}

func (b *badgerFSM) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:
		var payload = CommandPayload{}
		if err := json.Unmarshal(log.Data, &payload); err != nil {
			return nil
		}

		op := strings.ToUpper(strings.TrimSpace(payload.Operation))

		switch op {
		case "SET":
		case "GET":
		case "DELETE":
		}
	}

	return nil
}

func (b *badgerFSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (b *badgerFSM) Restore(io.ReadCloser) error {
	return nil
}
func (b *badgerFSM) set(key string, value interface{}) error {
	var data = make([]byte, 0)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if data == nil || len(data) <= 0 {
		return nil
	}

	txn := b.db.NewTransaction(true)
	if err := txn.Set([]byte(key), data); err != nil {
		return err
	}

	return txn.Commit()
}

func (b *badgerFSM) get(key string) (interface{}, error) {
	var keyByte = []byte(key)
	var data interface{}

	txn := b.db.NewTransaction(false)
	defer func() {
		txn.Commit()
	}()

	item, err := txn.Get(keyByte)
	if err != nil {
		data = map[string]interface{}{}
		return data, err
	}

	var value = make([]byte, 0)
	err = item.Value(func(val []byte) error {
		value = append(value, val...)
		return nil
	})

	if err != nil {
		data = map[string]interface{}{}
		return data, err
	}

	if len(value) > 0 {
		err = json.Unmarshal(value, &data)
	}

	if err != nil {
		data = map[string]interface{}{}
	}

	return data, err
}

func (b *badgerFSM) delete(key string) error {
	return nil
}

func NewBadgerFSM(db *badger.DB) *badgerFSM {
	return &badgerFSM{db}
}
