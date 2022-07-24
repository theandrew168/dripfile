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
				t, err := w.queue.Pop()
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
	// TODO: use a map
	// determine which task needs to run
	var handler TaskHandler
	switch t.Kind {
	case KindSessionPrune:
		handler = w.HandleSessionPrune
	case KindEmailSend:
		handler = w.HandleEmailSend
	case KindTransferTry:
		handler = w.HandleTransferTry
	default:
		w.logger.Error(fmt.Errorf("unknown task kind"), map[string]string{
			"task_id":   t.ID,
			"task_kind": string(t.Kind),
		})
		return
	}

	w.logger.Info("task start", map[string]string{
		"task_id":     t.ID,
		"task_kind":   string(t.Kind),
		"task_status": string(t.Status),
		"task_error":  t.Error,
	})

	ctx := context.Background()
	err := handler(ctx, t)
	if err != nil {
		w.logger.Error(err, map[string]string{
			"task_id":   t.ID,
			"task_kind": string(t.Kind),
		})

		// update status upon error
		t.Error = err.Error()
		t.Status = StatusError

		err = w.queue.Update(t)
		if err != nil {
			w.logger.Error(err, map[string]string{
				"task_id":   t.ID,
				"task_kind": string(t.Kind),
			})
			return
		}

		return
	}

	// delete the task upon success
	err = w.queue.Delete(t)
	if err != nil {
		w.logger.Error(err, map[string]string{
			"task_id":   t.ID,
			"task_kind": string(t.Kind),
		})
	}

	w.logger.Info("task finish", map[string]string{
		"task_id":     t.ID,
		"task_kind":   string(t.Kind),
		"task_status": string(t.Status),
		"task_error":  t.Error,
	})
}
