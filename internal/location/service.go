package location

type GetAllQuery struct{}

type GetByIDQuery struct {
	ID string
}

type AddMemoryCommand struct {
	ID string
}

type AddS3Command struct {
	ID string

	Endpoint        string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
}

type RemoveCommand struct {
	ID string
}

type Service interface {
	GetByID(query GetByIDQuery) (*Location, error)
	GetAll(query GetAllQuery) ([]*Location, error)
	AddMemory(cmd AddMemoryCommand) error
	AddS3(cmd AddS3Command) error
	Remove(cmd RemoveCommand) error
}
