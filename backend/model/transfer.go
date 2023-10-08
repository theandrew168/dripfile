package model

import "github.com/google/uuid"

type TransferStatus string

const (
	TransferStatusPending  = "pending"
	TransferStatusRunning  = "running"
	TransferStatusComplete = "complete"
)

type Transfer struct {
	ID uuid.UUID

	ItineraryID uuid.UUID
	Status      TransferStatus
	Progress    int
}

func NewTransfer(itineraryID uuid.UUID) Transfer {
	transfer := Transfer{
		ID: uuid.New(),

		ItineraryID: itineraryID,
		Status:      TransferStatusPending,
		Progress:    0,
	}
	return transfer
}

func (transfer Transfer) GetID() uuid.UUID {
	return transfer.ID
}
