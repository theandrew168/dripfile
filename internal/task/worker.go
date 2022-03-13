package task

import (
	"errors"
	"fmt"
	"time"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/log"
)

type Worker struct {
	queue   Queue
	storage database.Storage
	logger  log.Logger
}

func NewWorker(queue Queue, storage database.Storage, logger log.Logger) *Worker {
	worker := Worker{
		queue:   queue,
		storage: storage,
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
	w.logger.Info("task %s start\n", task.ID)
	switch task.Kind {
	case KindEmail:
	case KindSession:
		err := w.storage.Session.DeleteExpired()
		if err != nil {
			w.TaskError(task, err)
		}

		w.TaskSuccess(task)
	case KindTransfer:
	default:
		err := fmt.Errorf("unknown task: %s", task.Kind)
		w.logger.Error(err)
	}
	w.logger.Info("task %s finish\n", task.ID)
}

func (w *Worker) TaskSuccess(task Task) {
	task.Status = StatusSuccess

	err := w.queue.Update(task)
	if err != nil {
		w.logger.Error(err)
	}
}

func (w *Worker) TaskError(task Task, err error) {
	w.logger.Error(err)

	task.Status = StatusError
	task.Error = err.Error()

	err = w.queue.Update(task)
	if err != nil {
		w.logger.Error(err)
	}
}
