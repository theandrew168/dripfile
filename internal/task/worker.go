package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/log"
	"github.com/theandrew168/dripfile/internal/mail"
)

type TaskFunc func(task Task) error

type Worker struct {
	queue   Queue
	storage database.Storage
	mailer  mail.Mailer
	logger  log.Logger
}

func NewWorker(queue Queue, storage database.Storage, mailer mail.Mailer, logger log.Logger) *Worker {
	worker := Worker{
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
			task, err := w.queue.Pop()
			if err != nil {
				// break loop if no new tasks remain
				if errors.Is(err, core.ErrNotExist) {
					break
				}
				return err
			}

			// run task in the background
			go w.RunTask(task)
		}
	}

	return nil
}

func (w *Worker) RunTask(task Task) {
	// determine which task needs to run
	var taskFunc TaskFunc
	switch task.Kind {
	case KindEmail:
		taskFunc = w.SendEmail
	case KindSession:
		taskFunc = w.DeleteExpiredSessions
	case KindTransfer:
		taskFunc = w.DoTransfer
	default:
		err := fmt.Errorf("unknown task: %s", task.Kind)
		w.logger.Error(err)
		return
	}

	w.logger.Info("task %s start\n", task.ID)

	// run and update the status
	err := taskFunc(task)
	if err != nil {
		w.logger.Error(err)

		task.Error = err.Error()
		task.Status = StatusError
	} else {
		task.Status = StatusSuccess
	}

	err = w.queue.Update(task)
	if err != nil {
		w.logger.Error(err)
	}

	w.logger.Info("task %s finish\n", task.ID)
}

func (w *Worker) SendEmail(task Task) error {
	var info SendEmailInfo
	err := json.Unmarshal([]byte(task.Info), &info)
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

func (w *Worker) DeleteExpiredSessions(task Task) error {
	err := w.storage.Session.DeleteExpired()
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) DoTransfer(task Task) error {
	return nil
}
