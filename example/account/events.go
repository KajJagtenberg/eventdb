package account

type AccountEvent interface {
	Apply(acc *Account) error
}

type AccountRegistered struct {
	Name string `json:"name"`
}

func (e *AccountRegistered) Apply(acc *Account) error {
	acc.Name = e.Name

	return nil
}

type AccountNameChanged struct {
	Name string `json:"name"`
}

func (e *AccountNameChanged) Apply(acc *Account) error {
	acc.Name = e.Name

	return nil
}
