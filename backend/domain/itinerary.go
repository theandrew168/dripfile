package domain

import (
	"errors"
	"time"

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

	createdAt time.Time
	version   int
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

// Create an itinerary from existing data
// TODO: Can this be made visible ONLY to the repository package?
func LoadItinerary(id uuid.UUID, pattern string, fromLocationID, toLocationID uuid.UUID, createdAt time.Time, version int) *Itinerary {
	i := Itinerary{
		id: id,

		pattern:        pattern,
		fromLocationID: fromLocationID,
		toLocationID:   toLocationID,

		createdAt: createdAt,
		version:   version,
	}
	return &i
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

func (i *Itinerary) CreatedAt() time.Time {
	return i.createdAt
}

func (i *Itinerary) Version() int {
	return i.version
}

func (i *Itinerary) CheckDelete() error {
	return nil
}
