package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrItineraryInvalidPattern = errors.New("itinerary: invalid pattern")
	ErrItinerarySameLocation   = errors.New("itinerary: same location")
)

// Aggregate with a single entity
type Itinerary struct {
	id uuid.UUID

	pattern        string
	fromLocationID uuid.UUID
	toLocationID   uuid.UUID
}

// Factory func for creating a new itinerary
func NewItinerary(pattern string, from, to *Location) (*Itinerary, error) {
	if from.ID() == to.ID() {
		return nil, ErrItinerarySameLocation
	}
	if pattern == "" {
		return nil, ErrItineraryInvalidPattern
	}

	i := Itinerary{
		id: uuid.New(),

		pattern:        pattern,
		fromLocationID: from.ID(),
		toLocationID:   to.ID(),
	}

	from.useBy(&i)
	to.useBy(&i)

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

func (i *Itinerary) CheckDelete() error {
	return nil
}
