package location

import (
	"errors"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/location/fileserver"
	"github.com/theandrew168/dripfile/backend/location/fileserver/memory"
	"github.com/theandrew168/dripfile/backend/location/fileserver/s3"
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
	memoryInfo memory.Info
	s3Info     s3.Info
}

func NewMemory(id string) (*Location, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidUUID
	}

	info := memory.Info{}
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

	info := s3.Info{
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

func UnmarshalMemoryFromStorage(id string, info memory.Info) (*Location, error) {
	l := Location{
		id: id,

		kind:       KindMemory,
		memoryInfo: info,
	}
	return &l, nil
}

func UnmarshalS3FromStorage(id string, info s3.Info) (*Location, error) {
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

func (l *Location) MemoryInfo() memory.Info {
	return l.memoryInfo
}

func (l *Location) S3Info() s3.Info {
	return l.s3Info
}

func (l *Location) Connect() (fileserver.FileServer, error) {
	switch l.kind {
	case KindMemory:
		return memory.New(l.memoryInfo)
	case KindS3:
		return s3.New(l.s3Info)
	}

	return nil, ErrInvalidKind
}
