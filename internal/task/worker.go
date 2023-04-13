package task

import (
	"context"
	"errors"
	"time"

	"golang.org/x/exp/slog"
	"golang.org/x/sync/semaphore"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/mail"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/storage"
)

const (
	maxConcurrency = 1
)

type TaskHandler func(ctx context.Context, t Task) error

type Worker struct {
	logger *slog.Logger
	store  *storage.Storage
	queue  *Queue
	box    *secret.Box
	mailer mail.Mailer
}

func NewWorker(
	logger *slog.Logger,
	store *storage.Storage,
	queue *Queue,
	box *secret.Box,
	mailer mail.Mailer,
) *Worker {
	w := Worker{
		logger: logger,
		store:  store,
		queue:  queue,
		box:    box,
		mailer: mailer,
	}
	return &w
}

func (w *Worker) Run(ctx context.Context) error {
	sem := semaphore.NewWeighted(maxConcurrency)

	// check for new tasks periodically
	ticker := time.Tick(time.Second)
	for {
		select {
		case <-ctx.Done():
			goto stop
		case <-ticker:
			// kick off all new tasks
			for {
				// try to acquire semaphore slot,
				// loop again if concurrency is already maxed
				ok := sem.TryAcquire(1)
				if !ok {
					break
				}

				t, err := w.queue.Claim()
				if err != nil {
					// don't need the sem if no task was claimed
					sem.Release(1)

					// break loop if no new tasks remain
					if errors.Is(err, database.ErrNotExist) {
						break
					}
					return err
				}

				// run task in the background
				go func() {
					defer sem.Release(1)
					w.handleTask(t)
				}()
			}
		}
	}

stop:
	w.logger.Info("stopping worker")

	// wait for all tasks to finish
	sem.Acquire(context.Background(), maxConcurrency)

	w.logger.Info("stopped worker")
	return nil
}

func (w *Worker) handleTask(t Task) {
	// determine which task needs to run
	handlers := map[Kind]TaskHandler{
		KindSessionPrune: w.HandleSessionPrune,
		KindEmailSend:    w.HandleEmailSend,
	}

	handler, ok := handlers[t.Kind]
	if !ok {
		w.logger.Error("unknown task kind",
			slog.String("task_id", t.ID),
			slog.String("task_kind", string(t.Kind)),
		)
		return
	}

	w.logger.Info("task start",
		slog.String("task_id", t.ID),
		slog.String("task_kind", string(t.Kind)),
	)

	err := handler(context.Background(), t)
	if err != nil {
		// update status upon error
		t.Status = StatusFailure
		t.Error = err.Error()

		w.logger.Info("task failure",
			slog.String("task_id", t.ID),
			slog.String("task_kind", string(t.Kind)),
			slog.String("task_error", t.Error),
		)
	} else {
		w.logger.Info("task success",
			slog.String("task_id", t.ID),
			slog.String("task_kind", string(t.Kind)),
		)
	}

	err = w.queue.Finish(t)
	if err != nil {
		w.logger.Error(err.Error(),
			slog.String("task_id", t.ID),
			slog.String("task_kind", string(t.Kind)),
		)
		return
	}
}
