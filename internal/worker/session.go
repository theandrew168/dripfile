package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

func (w *Worker) HandleSessionPrune(ctx context.Context, t *asynq.Task) error {
	err := w.store.Session.DeleteExpired()
	if err != nil {
		return err
	}

	return nil
}
