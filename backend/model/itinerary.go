package model

import "github.com/google/uuid"

type Itinerary struct {
	ID uuid.UUID

	Pattern        string
	FromLocationID uuid.UUID
	ToLocationID   uuid.UUID
}

func NewItinerary(pattern string, fromLocation, toLocation Location) Itinerary {
	i := Itinerary{
		ID: uuid.New(),

		Pattern:        pattern,
		FromLocationID: fromLocation.ID,
		ToLocationID:   toLocation.ID,
	}
	return i
}

func (i Itinerary) GetID() uuid.UUID {
	return i.ID
}
