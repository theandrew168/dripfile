package domain

import (
	"time"

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

	createdAt time.Time
	updatedAt time.Time
}

func NewTransfer(itinerary *Itinerary) (*Transfer, error) {
	transfer := Transfer{
		id: uuid.New(),

		itineraryID: itinerary.ID(),
		status:      TransferStatusPending,
		progress:    0,
		error:       "",

		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
	return &transfer, nil
}

// Create an transfer from existing data
func LoadTransfer(
	id uuid.UUID,
	itineraryID uuid.UUID,
	status TransferStatus,
	progress int,
	error string,
	createdAt time.Time,
	updatedAt time.Time,
) *Transfer {
	t := Transfer{
		id: id,

		itineraryID: itineraryID,
		status:      status,
		progress:    progress,
		error:       error,

		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	return &t
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

func (t *Transfer) Error() string {
	return t.error
}

func (t *Transfer) UpdateError(error string) error {
	t.error = error
	return nil
}

func (t *Transfer) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Transfer) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Transfer) SetUpdatedAt(updatedAt time.Time) {
	t.updatedAt = updatedAt
}

func (t *Transfer) CheckDelete() error {
	return nil
}
