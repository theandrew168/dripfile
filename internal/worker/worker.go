package worker

import (
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/log"
	"github.com/theandrew168/dripfile/internal/pubsub"
)

type Worker struct {
	storage database.Storage
	queue   pubsub.Queue
	logger  log.Logger
}

func New(storage database.Storage, queue pubsub.Queue, logger log.Logger) *Worker {
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
	transfer, err := w.queue.Transfer.Subscribe()
	if err != nil {
		return err
	}

	w.logger.Info("transfer %s start\n", transfer.ID)
	w.logger.Info("transfer %s end\n", transfer.ID)
	return nil
}
