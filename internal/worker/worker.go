package worker

import (
	"errors"
	"time"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/log"
	"github.com/theandrew168/dripfile/internal/task"
)

type Worker struct {
	storage database.Storage
	queue   task.Queue
	logger  log.Logger
}

func New(storage database.Storage, queue task.Queue, logger log.Logger) *Worker {
	worker := Worker{
		storage: storage,
		queue:   queue,
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
	w.logger.Info("task %s start\n", t.ID)
	w.logger.Info("task %s end\n", t.ID)
}
