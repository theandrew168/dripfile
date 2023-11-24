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

	fromLocationID uuid.UUID
	toLocationID   uuid.UUID
	pattern        string

	createdAt time.Time
	updatedAt time.Time
}

// Factory func for creating a new itinerary
func NewItinerary(from, to *Location, pattern string) (*Itinerary, error) {
	if from.ID() == to.ID() {
		return nil, ErrItinerarySameLocation
	}
	if pattern == "" {
		return nil, ErrItineraryInvalidPattern
	}

	i := Itinerary{
		id: uuid.New(),

		fromLocationID: from.ID(),
		toLocationID:   to.ID(),
		pattern:        pattern,

		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	from.useBy(&i)
	to.useBy(&i)

	return &i, nil
}

// Create an itinerary from existing data
func LoadItinerary(
	id uuid.UUID,
	fromLocationID uuid.UUID,
	toLocationID uuid.UUID,
	pattern string,
	createdAt time.Time,
	updatedAt time.Time,
) *Itinerary {
	i := Itinerary{
		id: id,

		fromLocationID: fromLocationID,
		toLocationID:   toLocationID,
		pattern:        pattern,

		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	return &i
}

func (i *Itinerary) ID() uuid.UUID {
	return i.id
}

func (i *Itinerary) FromLocationID() uuid.UUID {
	return i.fromLocationID
}

func (i *Itinerary) ToLocationID() uuid.UUID {
	return i.toLocationID
}

func (i *Itinerary) Pattern() string {
	return i.pattern
}

func (i *Itinerary) CreatedAt() time.Time {
	return i.createdAt
}

func (i *Itinerary) UpdatedAt() time.Time {
	return i.updatedAt
}

func (i *Itinerary) SetUpdatedAt(updatedAt time.Time) {
	i.updatedAt = updatedAt
}

func (i *Itinerary) CheckDelete() error {
	return nil
}
