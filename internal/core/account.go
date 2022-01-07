package core

import (
	"context"
)

// enum values for account role
const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleViewer = "viewer"
)

// TODO: how does this look when using OAuth? via GitHub for example
type Account struct {
	Email    string
	Password string
	Verified bool
	Role     string
	Project  Project

	// readonly (from database, after creation)
	ID int
}

func NewAccount(email, password string, project Project) Account {
	account := Account{
		Email:    email,
		Password: password,
		Verified: false,
		Role:     RoleViewer,
		Project:  project,
	}
	return account
}

type AccountStorage interface {
	CreateAccount(ctx context.Context, account *Account) error
	ReadAccount(ctx context.Context, id int) (Account, error)
	UpdateAccount(ctx context.Context, account Account) error
	DeleteAccount(ctx context.Context, account Account) error
}
