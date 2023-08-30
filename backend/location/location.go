package location

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/fileserver"
)

const (
	KindMemory = "memory"
	KindS3     = "s3"
)

var (
	ErrInvalidUUID = errors.New("location: invalid UUID")
	ErrInvalidKind = errors.New("location: invalid kind")
)

type Location struct {
	id string

	kind       string
	memoryInfo fileserver.MemoryInfo
	s3Info     fileserver.S3Info

	// internal fields used for storage conflict resolution
	createdAt time.Time
	version   int
}

func NewMemory(id string) (*Location, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidUUID
	}

	info := fileserver.MemoryInfo{}
	err = info.Validate()
	if err != nil {
		return nil, err
	}

	l := Location{
		id: id,

		kind:       KindMemory,
		memoryInfo: info,
	}
	return &l, nil
}

func NewS3(id, endpoint, bucket, accessKeyID, secretAccessKey string) (*Location, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidUUID
	}

	info := fileserver.S3Info{
		Endpoint:        endpoint,
		Bucket:          bucket,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
	err = info.Validate()
	if err != nil {
		return nil, err
	}

	l := Location{
		id: id,

		kind:   KindS3,
		s3Info: info,
	}
	return &l, nil
}

func (l *Location) ID() string {
	return l.id
}

func (l *Location) Kind() string {
	return l.kind
}

func (l *Location) MemoryInfo() fileserver.MemoryInfo {
	return l.memoryInfo
}

func (l *Location) S3Info() fileserver.S3Info {
	return l.s3Info
}

func (l *Location) SetMemory(info fileserver.MemoryInfo) error {
	err := info.Validate()
	if err != nil {
		return err
	}

	l.kind = KindMemory
	l.memoryInfo = info
	return nil
}

func (l *Location) SetS3(info fileserver.S3Info) error {
	err := info.Validate()
	if err != nil {
		return err
	}

	l.kind = KindS3
	l.s3Info = info
	return nil
}

func (l *Location) Connect() (fileserver.FileServer, error) {
	switch l.kind {
	case KindMemory:
		return fileserver.NewMemory(l.memoryInfo)
	case KindS3:
		return fileserver.NewS3(l.s3Info)
	}

	return nil, ErrInvalidKind
}
