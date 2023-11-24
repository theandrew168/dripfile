package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/theandrew168/dripfile/backend/repository"
)

// References:
// https://brandur.org/postgres-queues
// https://webapp.io/blog/postgres-is-the-answer/
// https://www.2ndquadrant.com/en/blog/what-is-select-skip-locked-for-in-postgresql-9-5/

type Worker struct {
	logger *slog.Logger
	repo   *repository.Repository

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

	// run til ctx is cancelled
	// check for transfers in "pending" state
	// select w/ for update skip locked

	// sync.WaitGroup for running tasks (upper limit via semaphore?)
	// for+select loop w/ ticker (5s or something?), ctx.Done()
	// when done, break and wait for WG to finish

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	running := true
	for running {
		select {
		case <-ticker.C:
			w.logger.Info("checking for jobs")
			w.wg.Add(1)
			go func() {
				defer w.wg.Done()
				time.Sleep(3 * time.Second)
			}()
		case <-ctx.Done():
			running = false
		}
	}

	w.logger.Info("stopping worker")
	w.wg.Wait()

	w.logger.Info("stopped worker")

	return nil
}
