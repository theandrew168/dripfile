package transfer

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidPattern = errors.New("transfer: invalid pattern")
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

func New(id uuid.UUID, pattern string, fromLocationID, toLocationID uuid.UUID) (*Transfer, error) {
	if pattern == "" {
		return nil, ErrInvalidPattern
	}

	t := Transfer{
		id: id,

		pattern:        pattern,
		fromLocationID: fromLocationID,
		toLocationID:   toLocationID,
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
