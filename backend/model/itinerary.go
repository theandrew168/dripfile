package model

import "github.com/google/uuid"

type Itinerary struct {
	ID uuid.UUID

	Pattern        string
	FromLocationID uuid.UUID
	ToLocationID   uuid.UUID
}

func NewItinerary(pattern string, fromLocationID, toLocationID uuid.UUID) Itinerary {
	itinerary := Itinerary{
		ID: uuid.New(),

		Pattern:        pattern,
		FromLocationID: fromLocationID,
		ToLocationID:   toLocationID,
	}
	return itinerary
}

func (itinerary Itinerary) GetID() uuid.UUID {
	return itinerary.ID
}
