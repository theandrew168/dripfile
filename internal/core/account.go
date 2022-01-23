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
}

func NewAccount(email, username, password string) Account {
	account := Account{
		Email:    email,
		Username: username,
		Password: password,
		Verified: false,
		Role:     RoleViewer,
	}
	return account
}

type AccountStorage interface {
	// baseline CRUD ops all deal with one record
	Create(account *Account) error
	Read(id int64) (Account, error)
	Update(account Account) error
	Delete(account Account) error

	// ReadAll() ([]Account, error)
	// ReadManyByProject(project Project) ([]Account, error)

	// express the "many" side of one-to-many or many-to-many relationships

	// one account has many sessions, so:
	// SessionStorage: ReadManyByAccount(account Account) ([]Session, error)

	// many accounts belong to many projects, so:
	// AccountStorage: ReadManyByProject(project Project) ([]Account, error)
	// ProjectStorage: ReadManyByAccount(account Account) ([]Project, error)

	// TODO: flip this around to ByProject or something like that?
	Projects(account Account) ([]Project, error)
}
