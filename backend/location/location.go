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
	ErrInvalidUUID = errors.New("location: invalid UUID")
	ErrInvalidKind = errors.New("location: invalid kind")
)

type Location struct {
	id string

	kind       string
	memoryInfo fileserver.MemoryInfo
	s3Info     fileserver.S3Info
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

func (l *Location) Connect() (fileserver.FileServer, error) {
	switch l.kind {
	case KindMemory:
		return fileserver.NewMemory(l.memoryInfo)
	case KindS3:
		return fileserver.NewS3(l.s3Info)
	}

	return nil, ErrInvalidKind
}
