package work

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/log"
	"github.com/theandrew168/dripfile/internal/mail"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/task"
)

type TaskFunc func(t task.Task) error

type Worker struct {
	box     secret.Box
	queue   task.Queue
	storage database.Storage
	mailer  mail.Mailer
	logger  log.Logger
}

func NewWorker(box secret.Box, queue task.Queue, storage database.Storage, mailer mail.Mailer, logger log.Logger) *Worker {
	worker := Worker{
		box:     box,
		queue:   queue,
		storage: storage,
		mailer:  mailer,
		logger:  logger,
	}
	return &worker
}

// listen on queue, grab jobs, do the work, update as needed, success or error
func (w *Worker) Run() error {
	// check for new tasks periodically
	c := time.Tick(time.Second)
	for range c {
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
	}

	return nil
}

func (w *Worker) RunTask(t task.Task) {
	// determine which task needs to run
	var taskFunc TaskFunc
	switch t.Kind {
	case task.KindEmail:
		taskFunc = w.SendEmail
	case task.KindSession:
		taskFunc = w.DeleteExpiredSessions
	case task.KindTransfer:
		taskFunc = w.DoTransfer
	default:
		err := fmt.Errorf("unknown task: %s", t.Kind)
		w.logger.Error(err)
		return
	}

	w.logger.Info("task %s start\n", t.ID)

	// run and update the status
	err := taskFunc(t)
	if err != nil {
		w.logger.Error(err)

		t.Error = err.Error()
		t.Status = task.StatusError
	} else {
		t.Status = task.StatusSuccess
	}

	err = w.queue.Update(t)
	if err != nil {
		w.logger.Error(err)
	}

	w.logger.Info("task %s finish\n", t.ID)
}

func (w *Worker) SendEmail(t task.Task) error {
	var info task.SendEmailInfo
	err := json.Unmarshal([]byte(t.Info), &info)
	if err != nil {
		return err
	}

	err = w.mailer.SendEmail(
		info.FromName,
		info.FromEmail,
		info.ToName,
		info.ToEmail,
		info.Subject,
		info.Body,
	)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) DeleteExpiredSessions(t task.Task) error {
	err := w.storage.Session.DeleteExpired()
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) DoTransfer(t task.Task) error {
	return nil
}
