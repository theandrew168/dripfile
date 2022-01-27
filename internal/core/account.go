package core

// TODO: how does this look when using OAuth? via GitHub for example
type Account struct {
	// readonly (from database, after creation)
	ID string

	Email    string
	Username string
	Password string
	Verified bool
}

func NewAccount(email, username, password string) Account {
	account := Account{
		Email:    email,
		Username: username,
		Password: password,
		Verified: false,
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
	ReadManyByProject(project Project) ([]Account, error)
}
