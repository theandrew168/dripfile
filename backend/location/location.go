package location

import (
	"errors"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/fileserver"
)

const (
	KindMemory = "memory"
	KindS3     = "s3"
)

var (
	ErrInvalidKind = errors.New("location: invalid kind")
)

// Aggregate with a single entity
type Location struct {
	id uuid.UUID

	kind       string
	memoryInfo fileserver.MemoryInfo
	s3Info     fileserver.S3Info
}

// Factory func for creating a new in-memory location
func NewMemory() (*Location, error) {
	info := fileserver.MemoryInfo{}
	err := info.Validate()
	if err != nil {
		return nil, err
	}

	l := Location{
		id: uuid.New(),

		kind:       KindMemory,
		memoryInfo: info,
	}
	return &l, nil
}

// Factory func for creating a new S3 location
func NewS3(endpoint, bucket, accessKeyID, secretAccessKey string) (*Location, error) {
	info := fileserver.S3Info{
		Endpoint:        endpoint,
		Bucket:          bucket,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
	err := info.Validate()
	if err != nil {
		return nil, err
	}

	l := Location{
		id: uuid.New(),

		kind:   KindS3,
		s3Info: info,
	}
	return &l, nil
}

func (l *Location) ID() uuid.UUID {
	return l.id
}

func (l *Location) Kind() string {
	return l.kind
}

func (l *Location) Connect() (fileserver.FileServer, error) {
	switch l.kind {
	case KindMemory:
		return fileserver.NewMemory(l.memoryInfo)
	case KindS3:
		return fileserver.NewS3(l.s3Info)
	default:
		return nil, ErrInvalidKind
	}
}
