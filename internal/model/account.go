package model

// enum values for account / project relationship
const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleViewer = "viewer"
)

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
