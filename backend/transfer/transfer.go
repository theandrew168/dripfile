package transfer

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidUUID    = errors.New("transfer: invalid UUID")
	ErrInvalidPattern = errors.New("transfer: invalid pattern")
)

type Transfer struct {
	id string

	pattern        string
	fromLocationID string
	toLocationID   string
}

func New(id, pattern, fromLocationID, toLocationID string) (*Transfer, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidUUID
	}

	if pattern == "" {
		return nil, ErrInvalidPattern
	}

	_, err = uuid.Parse(fromLocationID)
	if err != nil {
		return nil, ErrInvalidUUID
	}

	_, err = uuid.Parse(toLocationID)
	if err != nil {
		return nil, ErrInvalidUUID
	}

	t := Transfer{
		id: id,

		pattern:        pattern,
		fromLocationID: fromLocationID,
		toLocationID:   toLocationID,
	}
	return &t, nil
}

func (t *Transfer) ID() string {
	return t.id
}

func (t *Transfer) Pattern() string {
	return t.pattern
}

func (t *Transfer) FromLocationID() string {
	return t.fromLocationID
}

func (t *Transfer) ToLocationID() string {
	return t.toLocationID
}
