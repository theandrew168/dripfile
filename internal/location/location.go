package location

// enum values for location kind
const (
	KindS3 = "s3"
)

type Location struct {
	// readonly (from database, after creation)
	ID string

	Kind string
	Info []byte
}

func New(kind string, info []byte) Location {
	location := Location{
		Kind: kind,
		Info: info,
	}
	return location
}

type Repository interface {
	Create(location *Location) error
	Read(id string) (Location, error)
	List() ([]Location, error)
	Update(location Location) error
	Delete(id string) error
}
