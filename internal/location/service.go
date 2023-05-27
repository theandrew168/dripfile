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

type Service interface {
	AddMemory(cmd AddMemoryCommand) error
	AddS3(cmd AddS3Command) error
	GetByID(query GetByIDQuery) (*Location, error)
	GetAll(query GetAllQuery) ([]*Location, error)
}
