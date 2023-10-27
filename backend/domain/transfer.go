package domain

import (
	"errors"

	"github.com/google/uuid"
)

type TransferStatus string

const (
	TransferStatusPending TransferStatus = "pending"
	TransferStatusRunning TransferStatus = "running"
	TransferStatusSuccess TransferStatus = "success"
	TransferStatusFailure TransferStatus = "failure"
)

type Transfer struct {
	id uuid.UUID

	itineraryID uuid.UUID
	status      TransferStatus
	progress    int
	error       string
}

func NewTransfer(itinerary *Itinerary) (*Transfer, error) {
	transfer := Transfer{
		id: uuid.New(),

		itineraryID: itinerary.ID(),
		status:      TransferStatusPending,
		progress:    0,
		error:       "",
	}
	return &transfer, nil
}

func (t *Transfer) ID() uuid.UUID {
	return t.id
}

func (t *Transfer) ItineraryID() uuid.UUID {
	return t.itineraryID
}

func (t *Transfer) Status() TransferStatus {
	return t.status
}

func (t *Transfer) UpdateStatus(status TransferStatus) error {
	t.status = status
	return nil
}

func (t *Transfer) Progress() int {
	return t.progress
}

func (t *Transfer) UpdateProgress(progress int) error {
	t.progress = progress
	return nil
}

func (t *Transfer) Error() error {
	if t.error != "" {
		return errors.New(t.error)
	}
	return nil
}

func (t *Transfer) UpdateError(error string) error {
	t.error = error
	return nil
}

func (t *Transfer) CheckDelete() error {
	return nil
}
