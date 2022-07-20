package task

import (
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
