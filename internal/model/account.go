package model

const (
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
}

func NewAccount(email, password, role string) Account {
	account := Account{
		Email:    email,
		Password: password,
		Role:     role,
		Verified: false,
	}
	return account
}
