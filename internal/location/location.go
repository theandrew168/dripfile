package location

// enum values for location kind
const (
	KindS3 = "s3"
)

type Location struct {
	// readonly (from database, after creation)
	ID string

	Kind string
	Name string
	Info []byte
}

func New(kind, name string, info []byte) Location {
	location := Location{
		Kind: kind,
		Name: name,
		Info: info,
	}
	return location
}
