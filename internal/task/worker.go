package task

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/mail"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/storage"
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
	var wg sync.WaitGroup

	// check for new tasks periodically
	ticker := time.Tick(time.Second)
	for {
		select {
		case <-w.stop:
			goto stop
		case <-ticker:
			// kick off all new tasks
			for {
				t, err := w.queue.Claim()
				if err != nil {
					// break loop if no new tasks remain
					if errors.Is(err, postgresql.ErrNotExist) {
						break
					}
					return err
				}

				// run task in the background
				wg.Add(1)
				go func() {
					defer wg.Done()
					w.handleTask(t)
				}()
			}
		}
	}

stop:
	wg.Wait()
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
