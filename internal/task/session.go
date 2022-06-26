package task

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeSessionPrune = "session:prune"
)

type SessionPrunePayload struct{}

func NewSessionPruneTask() (*asynq.Task, error) {
	payload := SessionPrunePayload{}

	js, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSessionPrune, js), nil
}

func (w *Worker) HandleSessionPrune(ctx context.Context, t *asynq.Task) error {
	err := w.store.Session.DeleteExpired()
	if err != nil {
		return err
	}

	return nil
}
