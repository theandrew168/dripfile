package worker

import (
	"context"

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
	for {
		// TODO: manual check every so often
		ctx := context.Background()

		// bail out if the listen fails
		err := w.queue.Listen(ctx)
		if err != nil {
			return err
		}

		// process a task
		task, err := w.queue.Pop()
		if err != nil {
			return err
		}

		w.logger.Info("task %s start\n", task.ID)
		w.logger.Info("task %s end\n", task.ID)
	}

	return nil
}
