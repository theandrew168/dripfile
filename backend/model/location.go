package model

import (
	"errors"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/fileserver"
)

const (
	LocationKindMemory = "memory"
	LocationKindS3     = "s3"
)

var (
	ErrInvalidLocationKind = errors.New("location: invalid kind")
)

type Location struct {
	ID uuid.UUID

	Kind       string
	MemoryInfo fileserver.MemoryInfo
	S3Info     fileserver.S3Info
}

func NewMemoryLocation() Location {
	info := fileserver.MemoryInfo{}
	l := Location{
		ID: uuid.New(),

		Kind:       LocationKindMemory,
		MemoryInfo: info,
	}
	return l
}

func NewS3Location(endpoint, bucket, accessKeyID, secretAccessKey string) Location {
	info := fileserver.S3Info{
		Endpoint:        endpoint,
		Bucket:          bucket,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
	l := Location{
		ID: uuid.New(),

		Kind:   LocationKindS3,
		S3Info: info,
	}
	return l
}

func (l Location) GetID() uuid.UUID {
	return l.ID
}

func (l Location) Connect() (fileserver.FileServer, error) {
	switch l.Kind {
	case LocationKindMemory:
		return fileserver.NewMemory(l.MemoryInfo)
	case LocationKindS3:
		return fileserver.NewS3(l.S3Info)
	default:
		return nil, ErrInvalidLocationKind
	}
}
