package core

import (
	"context"
)

type Account struct {
	Name string

	// readonly (from database, after creation)
	ID int
}

func NewAccount(name string) Account {
	account := Account{
		Name: name,
	}
	return account
}

type AccountStorage interface {
	CreateAccount(ctx context.Context, account *Account) error
	ReadAccount(ctx context.Context, id int) (Account, error)
	UpdateAccount(ctx context.Context, account Account) error
	DeleteAccount(ctx context.Context, account Account) error
}
