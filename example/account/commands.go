package account

import (
	"encoding/json"
	"errors"

	"github.com/kajjagtenberg/eventflowdb/api"
)

type AccountCommand interface {
	Execute(acc *Account) ([]*api.EventData, error)
}

type RegisterAccount struct {
	Name string `json:"name"`
}

func (cmd *RegisterAccount) Execute(acc *Account) ([]*api.EventData, error) {
	if acc.Version > 0 {
		return nil, errors.New("account already exists")
	}

	data, err := json.Marshal(&AccountRegistered{
		Name: cmd.Name,
	})
	if err != nil {
		return nil, err
	}

	return []*api.EventData{
		{
			Type: "AccountRegistered",
			Data: data,
		},
	}, nil
}

type ChangeAccountName struct {
	Name string `json:"name"`
}

func (cmd *ChangeAccountName) Execute(acc *Account) ([]*api.EventData, error) {
	if acc.Version == 0 {
		return nil, errors.New("account does not exist")
	}

	data, err := json.Marshal(&AccountNameChanged{
		Name: cmd.Name,
	})
	if err != nil {
		return nil, err
	}

	return []*api.EventData{
		{
			Type: "AccountNameChanged",
			Data: data,
		},
	}, nil
}
