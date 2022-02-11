package core

// enum values for location kind
const (
	KindS3   = "s3"
	KindFTP  = "ftp"
	KindFTPS = "ftps"
	KindSFTP = "sftp"
)

type Location struct {
	// readonly (from database, after creation)
	ID string

	Kind    string
	Info    string
	Project Project
}

func NewLocation(kind, info string, project Project) Location {
	location := Location{
		Kind:    kind,
		Info:    info,
		Project: project,
	}
	return location
}

type LocationStorage interface {
	Create(location *Location) error
	Read(id string) (Location, error)
	Update(location Location) error
	Delete(location Location) error

	ReadManyByProject(project Project) ([]Location, error)
}
