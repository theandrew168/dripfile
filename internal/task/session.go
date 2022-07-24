package task

import (
	"context"
	"encoding/json"
)

const (
	KindSessionPrune = "session:prune"
)

type SessionPruneInfo struct{}

func NewSessionPruneTask() Task {
	info := SessionPruneInfo{}

	js, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	return NewTask(KindSessionPrune, string(js))
}

func (w *Worker) HandleSessionPrune(ctx context.Context, t Task) error {
	err := w.store.Session.DeleteExpired()
	if err != nil {
		return err
	}

	return nil
}
