package core

// enum values for account role
const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleViewer = "viewer"
)

// TODO: how does this look when using OAuth? via GitHub for example
type Account struct {
	// readonly (from database, after creation)
	ID int64

	Email    string
	Username string
	Password string
	Verified bool
	Role     string
	Project  Project
}

func NewAccount(email, username, password string, project Project) Account {
	account := Account{
		Email:    email,
		Username: username,
		Password: password,
		Verified: false,
		Role:     RoleViewer,
		Project:  project,
	}
	return account
}

type AccountStorage interface {
	Create(account *Account) error
	Read(id int64) (Account, error)
	Update(account Account) error
	Delete(account Account) error
}
