package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/fileserver"
)

type LocationKind string

const (
	LocationKindMemory LocationKind = "memory"
	LocationKindS3     LocationKind = "s3"
)

type PingStatus string

const (
	PingStatusUnknown PingStatus = "unknown"
	PingStatusSuccess PingStatus = "success"
	PingStatusFailure PingStatus = "failure"
)

var (
	ErrLocationInvalidKind = errors.New("location: invalid kind")

	// TODO: In use by what?
	ErrLocationInUse = errors.New("location: in use")
)

// Aggregate with a single entity
type Location struct {
	id uuid.UUID

	kind       LocationKind
	memoryInfo fileserver.MemoryInfo
	s3Info     fileserver.S3Info
	pingStatus PingStatus

	createdAt time.Time
	updatedAt time.Time

	usedBy []uuid.UUID
}

// Factory func for creating a new in-memory location
func NewMemoryLocation() (*Location, error) {
	info := fileserver.MemoryInfo{}

	l := Location{
		id: uuid.New(),

		kind:       LocationKindMemory,
		memoryInfo: info,
		pingStatus: PingStatusUnknown,

		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
	return &l, nil
}

// Create an in-memory location from existing data
func LoadMemoryLocation(
	id uuid.UUID,
	info fileserver.MemoryInfo,
	pingStatus PingStatus,
	createdAt time.Time,
	updatedAt time.Time,
	usedBy []uuid.UUID,
) *Location {
	l := Location{
		id: id,

		kind:       LocationKindMemory,
		memoryInfo: info,
		pingStatus: pingStatus,

		createdAt: createdAt,
		updatedAt: updatedAt,

		usedBy: usedBy,
	}
	return &l
}

// Factory func for creating a new S3 location
func NewS3Location(endpoint, bucket, accessKeyID, secretAccessKey string) (*Location, error) {
	if endpoint == "" {
		return nil, errors.New("location: empty S3 endpoint")
	}
	if bucket == "" {
		return nil, errors.New("location: empty S3 bucket")
	}
	if accessKeyID == "" {
		return nil, errors.New("location: empty S3 access key id")
	}
	if secretAccessKey == "" {
		return nil, errors.New("location: empty S3 secret access key")
	}

	info := fileserver.S3Info{
		Endpoint:        endpoint,
		Bucket:          bucket,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}

	l := Location{
		id: uuid.New(),

		kind:       LocationKindS3,
		s3Info:     info,
		pingStatus: PingStatusUnknown,

		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
	return &l, nil
}

// Create an S3 location from existing data
func LoadS3Location(
	id uuid.UUID,
	info fileserver.S3Info,
	pingStatus PingStatus,
	createdAt time.Time,
	updatedAt time.Time,
	usedBy []uuid.UUID,
) *Location {
	l := Location{
		id: id,

		kind:       LocationKindS3,
		s3Info:     info,
		pingStatus: pingStatus,

		createdAt: createdAt,
		updatedAt: updatedAt,

		usedBy: usedBy,
	}
	return &l
}

func (l *Location) ID() uuid.UUID {
	return l.id
}

func (l *Location) Kind() LocationKind {
	return l.kind
}

func (l *Location) Info() any {
	if l.kind == LocationKindMemory {
		return l.memoryInfo
	} else if l.kind == LocationKindS3 {
		return l.s3Info
	}

	return nil
}

func (l *Location) PingStatus() PingStatus {
	return l.pingStatus
}

func (l *Location) SetPingStatus(pingStatus PingStatus) error {
	l.pingStatus = pingStatus
	return nil
}

func (l *Location) CreatedAt() time.Time {
	return l.createdAt
}

func (l *Location) UpdatedAt() time.Time {
	return l.updatedAt
}

func (l *Location) SetUpdatedAt(updatedAt time.Time) error {
	l.updatedAt = updatedAt
	return nil
}

func (l *Location) UsedBy() []uuid.UUID {
	return l.usedBy
}

func (l *Location) CheckDelete() error {
	if len(l.usedBy) > 0 {
		return ErrLocationInUse
	}

	return nil
}

func (l *Location) Connect() (fileserver.FileServer, error) {
	switch l.kind {
	case LocationKindMemory:
		return fileserver.NewMemory(l.memoryInfo)
	case LocationKindS3:
		return fileserver.NewS3(l.s3Info)
	default:
		return nil, ErrLocationInvalidKind
	}
}

func (l *Location) useBy(itinerary *Itinerary) {
	l.usedBy = append(l.usedBy, itinerary.ID())
}
