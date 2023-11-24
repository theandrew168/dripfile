package worker

import (
	"context"

	"github.com/theandrew168/dripfile/backend/repository"
)

// References:
// https://brandur.org/postgres-queues
// https://webapp.io/blog/postgres-is-the-answer/
// https://www.2ndquadrant.com/en/blog/what-is-select-skip-locked-for-in-postgresql-9-5/

type Worker struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Worker {
	w := Worker{
		repo: repo,
	}
	return &w
}

func (w *Worker) Run(ctx *context.Context) error {
	// run til ctx is cancelled
	// check for transfers in "pending" state
	// select w/ for update skip locked

	// sync.WaitGroup for running tasks (upper limit via semaphore?)
	// for+select loop w/ ticker (5s or something?), ctx.Done()
	// when done, break and wait for WG to finish
	return nil
}
