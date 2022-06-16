package task

import (
	"errors"
	"fmt"
	"time"

	"github.com/theandrew168/dripfile/backend/core"
	"github.com/theandrew168/dripfile/backend/jsonlog"
	"github.com/theandrew168/dripfile/backend/mail"
	"github.com/theandrew168/dripfile/backend/secret"
	"github.com/theandrew168/dripfile/backend/storage"
	"github.com/theandrew168/dripfile/backend/stripe"
)

type TaskFunc func(task Task) error

type Worker struct {
	logger  *jsonlog.Logger
	storage *storage.Storage
	queue   *Queue
	box     *secret.Box
	billing stripe.Billing
	mailer  mail.Mailer
}

func NewWorker(
	logger *jsonlog.Logger,
	storage *storage.Storage,
	queue *Queue,
	box *secret.Box,
	billing stripe.Billing,
	mailer mail.Mailer,
) *Worker {
	worker := Worker{
		logger:  logger,
		storage: storage,
		queue:   queue,
		box:     box,
		billing: billing,
		mailer:  mailer,
	}
	return &worker
}

// listen on queue, grab jobs, do the work, update as needed, success or error
func (w *Worker) Run() error {
	// check for new tasks periodically
	c := time.Tick(time.Second)
	for {
		// kick off all new tasks
		for {
			t, err := w.queue.Pop()
			if err != nil {
				// break loop if no new tasks remain
				if errors.Is(err, core.ErrNotExist) {
					break
				}
				return err
			}

			// run task in the background
			go w.RunTask(t)
		}

		<-c
	}

	return nil
}

func (w *Worker) RunTask(task Task) {
	// determine which task needs to run
	var taskFunc TaskFunc
	switch task.Kind {
	case KindPruneSessions:
		taskFunc = w.PruneSessions
	case KindSendEmail:
		taskFunc = w.SendEmail
	case KindTransfer:
		taskFunc = w.Transfer
	default:
		w.logger.Error(fmt.Errorf("unknown task kind"), map[string]string{
			"task_id":   task.ID,
			"task_kind": task.Kind,
		})
		return
	}

	w.logger.Info("task start", map[string]string{
		"task_id":     task.ID,
		"task_kind":   task.Kind,
		"task_status": task.Status,
		"task_error":  task.Error,
	})

	// run and update the status
	err := taskFunc(task)
	if err != nil {
		w.logger.Error(err, map[string]string{
			"task_id":   task.ID,
			"task_kind": task.Kind,
		})

		task.Error = err.Error()
		task.Status = StatusError
	} else {
		task.Status = StatusSuccess
	}

	err = w.queue.Update(task)
	if err != nil {
		w.logger.Error(err, map[string]string{
			"task_id":   task.ID,
			"task_kind": task.Kind,
		})
	}

	w.logger.Info("task finish", map[string]string{
		"task_id":     task.ID,
		"task_kind":   task.Kind,
		"task_status": task.Status,
		"task_error":  task.Error,
	})
}
