package core

import (
	"github.com/theandrew168/dripfile/pkg/random"
)

// enum values for account / project relationship
const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleViewer = "viewer"
)

// TODO: how does this look when using OAuth? via GitHub for example
type Account struct {
	// readonly (from database, after creation)
	ID string

	Email    string
	Password string
	Role     string
	Verified bool
	Project  Project
}

func NewAccount(email, password, role string, project Project) Account {
	account := Account{
		Email:    email,
		Password: password,
		Role:     role,
		Verified: false,
		Project:  project,
	}
	return account
}

func NewAccountMock(project Project) Account {
	account := NewAccount(
		random.String(8),
		random.String(8),
		RoleViewer,
		project,
	)
	return account
}
