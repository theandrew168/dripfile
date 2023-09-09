package itinerary

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidPattern = errors.New("itinerary: invalid pattern")
	ErrSameLocation   = errors.New("itinerary: same location")
)

type Itinerary struct {
	id uuid.UUID

	fromLocationID uuid.UUID
	toLocationID   uuid.UUID
	pattern        string
}

func New(fromLocationID, toLocationID uuid.UUID, pattern string) (*Itinerary, error) {
	if fromLocationID == toLocationID {
		return nil, ErrSameLocation
	}
	if pattern == "" {
		return nil, ErrInvalidPattern
	}

	i := Itinerary{
		id: uuid.New(),

		fromLocationID: fromLocationID,
		toLocationID:   toLocationID,
		pattern:        pattern,
	}
	return &i, nil
}

func (i *Itinerary) ID() uuid.UUID {
	return i.id
}

func (i *Itinerary) Pattern() string {
	return i.pattern
}

func (i *Itinerary) FromLocationID() uuid.UUID {
	return i.fromLocationID
}

func (i *Itinerary) ToLocationID() uuid.UUID {
	return i.toLocationID
}
