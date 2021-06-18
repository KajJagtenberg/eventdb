package account

import (
	"context"
	"encoding/json"

	"github.com/eventflowdb/eventflowdb/api"
	"github.com/google/uuid"
)

type Account struct {
	ID      uuid.UUID
	Version uint32

	Name string
}

func (acc *Account) LoadFromHistory(eventstore api.EventStoreClient) error {
	res, err := eventstore.Get(context.Background(), &api.GetRequest{
		Stream:  acc.ID.String(),
		Version: 0,
		Limit:   0,
	})
	if err != nil {
		return err
	}

	for _, event := range res.Events {
		if err := acc.Apply(event); err != nil {
			return err
		}
	}

	return nil
}

func (acc *Account) Apply(event *api.Event) error {
	var e AccountEvent

	switch event.Type {
	case "AccountRegistered":
		e = &AccountRegistered{}
	case "AccountNameChanged":
		e = &AccountNameChanged{}
	}

	if err := json.Unmarshal(event.Data, &e); err != nil {
		return err
	}

	if err := e.Apply(acc); err != nil {
		return nil
	}

	acc.Version++

	return nil
}

func Handle(eventstore api.EventStoreClient, id uuid.UUID, cmd AccountCommand) error {
	acc := &Account{
		ID: id,
	}

	if err := acc.LoadFromHistory(eventstore); err != nil {
		return err
	}

	events, err := cmd.Execute(acc)
	if err != nil {
		return err
	}

	_, err = eventstore.Add(context.Background(), &api.AddRequest{
		Stream:  id.String(),
		Version: acc.Version,
		Events:  events,
	})
	return err
}
