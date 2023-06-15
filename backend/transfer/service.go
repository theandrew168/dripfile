package transfer

type GetAllQuery struct{}

type GetByIDQuery struct {
	ID string
}

type AddCommand struct {
	ID string

	Pattern        string
	FromLocationID string
	ToLocationID   string
}

type RemoveCommand struct {
	ID string
}

type RunCommand struct {
	ID string
}

type Service interface {
	GetByID(query GetByIDQuery) (*Transfer, error)
	GetAll(query GetAllQuery) ([]*Transfer, error)

	Add(cmd AddCommand) error
	Remove(cmd RemoveCommand) error
	Run(cmd RunCommand) error
}
