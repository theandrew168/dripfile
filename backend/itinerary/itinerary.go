package itinerary

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidPattern = errors.New("itinerary: invalid pattern")
	ErrSameLocation   = errors.New("itinerary: same location")
)

// Aggregate with a single entity
type Itinerary struct {
	id uuid.UUID

	pattern        string
	fromLocationID uuid.UUID
	toLocationID   uuid.UUID
}

// Factory func for creating a new itinerary
func New(pattern string, fromLocationID, toLocationID uuid.UUID) (*Itinerary, error) {
	if fromLocationID == toLocationID {
		return nil, ErrSameLocation
	}
	if pattern == "" {
		return nil, ErrInvalidPattern
	}

	i := Itinerary{
		id: uuid.New(),

		pattern:        pattern,
		fromLocationID: fromLocationID,
		toLocationID:   toLocationID,
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
