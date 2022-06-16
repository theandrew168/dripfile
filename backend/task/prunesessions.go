package task

import (
	"encoding/json"
)

const KindPruneSessions = "prune_sessions"

type PruneSessionsInfo struct{}

func PruneSessions() (Task, error) {
	info := PruneSessionsInfo{}

	b, err := json.Marshal(info)
	if err != nil {
		return Task{}, err
	}

	return New(KindPruneSessions, string(b)), nil
}

func (w *Worker) PruneSessions(task Task) error {
	err := w.storage.Session.DeleteExpired()
	if err != nil {
		return err
	}

	return nil
}
