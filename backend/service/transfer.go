package service

import (
	"log/slog"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/fileserver"
	"github.com/theandrew168/dripfile/backend/repository"
)

type Transfer struct {
	logger *slog.Logger
	repo   *repository.Repository
}

func NewTransfer(logger *slog.Logger, repo *repository.Repository) *Transfer {
	t := Transfer{
		logger: logger,
		repo:   repo,
	}
	return &t
}

func (t *Transfer) Run(id uuid.UUID) error {
	// look up transfer by ID
	transfer, err := t.repo.Transfer.Read(id)
	if err != nil {
		return err
	}

	// look up itinerary by ID
	itinerary, err := t.repo.Itinerary.Read(transfer.ID())
	if err != nil {
		return err
	}

	// look up locations by ID
	fromLocation, err := t.repo.Location.Read(itinerary.FromLocationID())
	if err != nil {
		return err
	}

	toLocation, err := t.repo.Location.Read(itinerary.ToLocationID())
	if err != nil {
		return err
	}

	// connect to file servers
	from, err := fromLocation.Connect()
	if err != nil {
		return err
	}

	to, err := toLocation.Connect()
	if err != nil {
		return err
	}

	// update xfer status to in-progress
	err = transfer.UpdateStatus(domain.TransferStatusRunning)
	if err != nil {
		return err
	}

	err = t.repo.Transfer.Update(transfer)
	if err != nil {
		return err
	}

	// run the xfer
	// TODO: update the transfer (in DB) every N seconds
	progress, err := fileserver.Transfer(itinerary.Pattern(), from, to)
	if err != nil {
		return err
	}

	// update xfer status to success / failure + xfer progress
	err = transfer.UpdateStatus(domain.TransferStatusSuccess)
	if err != nil {
		return err
	}

	err = transfer.UpdateProgress(progress)
	if err != nil {
		return err
	}

	err = t.repo.Transfer.Update(transfer)
	if err != nil {
		return err
	}

	return nil
}
