package task

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/mail"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/storage"
)

const (
	maxConcurrency = 1
)

type TaskHandler func(ctx context.Context, t Task) error

type Worker struct {
	logger *jsonlog.Logger
	store  *storage.Storage
	queue  *Queue
	box    *secret.Box
	mailer mail.Mailer

	stop chan struct{}
}

func NewWorker(
	logger *jsonlog.Logger,
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

		stop: make(chan struct{}),
	}
	return &w
}

func (w *Worker) Start() error {
	sem := semaphore.NewWeighted(maxConcurrency)

	// check for new tasks periodically
	ticker := time.Tick(time.Second)
	for {
		select {
		case <-w.stop:
			goto stop
		case <-ticker:
			// kick off all new tasks
			for {
				// acquire semaphore slot
				sem.Acquire(context.Background(), 1)

				t, err := w.queue.Claim()
				if err != nil {
					// don't need the sem if no task was claimed
					sem.Release(1)

					// break loop if no new tasks remain
					if errors.Is(err, postgresql.ErrNotExist) {
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
	// wait for all tasks to finish
	sem.Acquire(context.Background(), maxConcurrency)
	return nil
}

func (w *Worker) Stop() error {
	w.logger.Info("stopping worker", nil)
	w.stop <- struct{}{}
	return nil
}

func (w *Worker) handleTask(t Task) {
	// determine which task needs to run
	handlers := map[Kind]TaskHandler{
		KindSessionPrune: w.HandleSessionPrune,
		KindEmailSend:    w.HandleEmailSend,
		KindTransferTry:  w.HandleTransferTry,
	}

	handler, ok := handlers[t.Kind]
	if !ok {
		w.logger.Error(fmt.Errorf("unknown task kind"), map[string]string{
			"task_id":   t.ID,
			"task_kind": string(t.Kind),
		})
		return
	}

	w.logger.Info("task start", map[string]string{
		"task_id":   t.ID,
		"task_kind": string(t.Kind),
	})

	err := handler(context.Background(), t)
	if err != nil {
		// update status upon error
		t.Status = StatusFailure
		t.Error = err.Error()

		w.logger.Info("task failure", map[string]string{
			"task_id":    t.ID,
			"task_kind":  string(t.Kind),
			"task_error": t.Error,
		})
	} else {
		w.logger.Info("task success", map[string]string{
			"task_id":   t.ID,
			"task_kind": string(t.Kind),
		})
	}

	err = w.queue.Finish(t)
	if err != nil {
		w.logger.Error(err, map[string]string{
			"task_id":   t.ID,
			"task_kind": string(t.Kind),
		})
		return
	}
}
