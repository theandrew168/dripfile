package core

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
	Username string
	Password string
	Role     string
	Verified bool
	Project  Project
}

func NewAccount(email, username, password, role string, project Project) Account {
	account := Account{
		Email:    email,
		Username: username,
		Password: password,
		Role:     role,
		Verified: false,
		Project:  project,
	}
	return account
}

type AccountStorage interface {
	// baseline CRUD ops all deal with one record
	Create(account *Account) error
	Read(id string) (Account, error)
	Update(account Account) error
	Delete(account Account) error

	// additional CRUD ops may deal with many
	ReadByEmail(email string) (Account, error)
}
