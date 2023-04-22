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
