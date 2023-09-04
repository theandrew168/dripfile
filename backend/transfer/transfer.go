package transfer

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/location"
)

var (
	ErrInvalidPattern = errors.New("transfer: invalid pattern")
	ErrSameLocation   = errors.New("transfer: same location")
)

type Transfer struct {
	id uuid.UUID

	pattern        string
	fromLocationID uuid.UUID
	toLocationID   uuid.UUID

	// internal fields used for storage conflict resolution
	createdAt time.Time
	version   int
}

func New(id uuid.UUID, pattern string, fromLocation, toLocation *location.Location) (*Transfer, error) {
	if pattern == "" {
		return nil, ErrInvalidPattern
	}
	if fromLocation.ID() == toLocation.ID() {
		return nil, ErrSameLocation
	}

	t := Transfer{
		id: id,

		pattern:        pattern,
		fromLocationID: fromLocation.ID(),
		toLocationID:   toLocation.ID(),
	}
	return &t, nil
}

func (t *Transfer) ID() uuid.UUID {
	return t.id
}

func (t *Transfer) Pattern() string {
	return t.pattern
}

func (t *Transfer) FromLocationID() uuid.UUID {
	return t.fromLocationID
}

func (t *Transfer) ToLocationID() uuid.UUID {
	return t.toLocationID
}
