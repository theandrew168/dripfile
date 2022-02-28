package worker

import (
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

func (w *Worker) Run() error {
	// listen on queue, grab jobs, do the work, update as needed, success or error

	// simulate a single job
	task, err := w.queue.Subscribe()
	if err != nil {
		return err
	}

	w.logger.Info("task %s start\n", task.ID)
	w.logger.Info("task %s end\n", task.ID)
	return nil
}
