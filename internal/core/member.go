package core

// enum values for account / project relationship
const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleViewer = "viewer"
)

// relationship between Accounts and Projects
type Member struct {
	// readonly (from database, after creation)
	ID string

	Role    string
	Account Account
	Project Project
}

func NewMember(role string, account Account, project Project) Member {
	member := Member{
		Role:    role,
		Account: account,
		Project: project,
	}
	return member
}

type MemberStorage interface {
	// baseline CRUD ops all deal with one record
	Create(member *Member) error
	Read(id string) (Member, error)
	Update(member Member) error
	Delete(member Member) error

	// additional CRUD ops may deal with many
	ReadManyByAccount(account Account) ([]Member, error)
	ReadManyByProject(project Project) ([]Member, error)
}
