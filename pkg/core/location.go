package core

import (
	"github.com/theandrew168/dripfile/pkg/random"
)

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
	Name    string
	Info    []byte
	Project Project
}

func NewLocation(kind, name string, info []byte, project Project) Location {
	location := Location{
		Kind:    kind,
		Name:    name,
		Info:    info,
		Project: project,
	}
	return location
}

func NewLocationMock(project Project) Location {
	location := NewLocation(
		random.String(8),
		random.String(8),
		random.Bytes(8),
		project,
	)
	return location
}
