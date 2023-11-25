package worker

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/fileserver"
	"github.com/theandrew168/dripfile/backend/repository"
)

// References:
// https://brandur.org/postgres-queues
// https://webapp.io/blog/postgres-is-the-answer/
// https://www.2ndquadrant.com/en/blog/what-is-select-skip-locked-for-in-postgresql-9-5/

type Worker struct {
	logger *slog.Logger
	repo   *repository.Repository

	// sync.WaitGroup for running tasks (upper limit via semaphore?)
	wg sync.WaitGroup
}

func New(logger *slog.Logger, repo *repository.Repository) *Worker {
	w := Worker{
		logger: logger,
		repo:   repo,
	}
	return &w
}

func (w *Worker) Run(ctx context.Context) error {
	w.logger.Info("starting worker")

	// do an initial poll before starting the ticker
	err := w.Poll()
	if err != nil {
		// log error but don't abort
		w.logger.Error(err.Error())
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	running := true
	for running {
		select {
		case <-ticker.C:
			err := w.Poll()
			if err != nil {
				// log error but don't abort
				w.logger.Error(err.Error())
			}
		case <-ctx.Done():
			running = false
		}
	}

	w.logger.Info("stopping worker")
	w.wg.Wait()

	w.logger.Info("stopped worker")

	return nil
}

func (w *Worker) Poll() error {
	w.logger.Info("checking for new jobs")

	// kick off all pending transfers
	for {
		transfer, err := w.repo.Transfer.Acquire()
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotExist):
				return nil
			default:
				return err
			}
		}

		w.logger.Info("running transfer", "id", transfer.ID())
		err = w.RunTransfer(transfer)
		if err != nil {
			transfer.UpdateError(err.Error())
			transfer.UpdateStatus(domain.TransferStatusFailure)
		} else {
			transfer.UpdateStatus(domain.TransferStatusSuccess)
		}

		err = w.repo.Transfer.Update(transfer)
		if err != nil {
			return err
		}
	}
}

func (w *Worker) RunTransfer(transfer *domain.Transfer) error {
	// look up itinerary by ID
	itinerary, err := w.repo.Itinerary.Read(transfer.ItineraryID())
	if err != nil {
		return err
	}

	// look up locations by ID
	fromLocation, err := w.repo.Location.Read(itinerary.FromLocationID())
	if err != nil {
		return err
	}

	toLocation, err := w.repo.Location.Read(itinerary.ToLocationID())
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

	// run the xfer
	// TODO: update the transfer (in DB) every N seconds
	progress, err := fileserver.Transfer(itinerary.Pattern(), from, to)
	if err != nil {
		return err
	}

	// TODO: update the xfer progress periodically
	err = transfer.UpdateProgress(progress)
	if err != nil {
		return err
	}

	err = w.repo.Transfer.Update(transfer)
	if err != nil {
		return err
	}

	return nil
}
