package core

import (
	"time"
)

type History struct {
	// readonly (from database, after creation)
	ID string

	Bytes      int64
	Status     string
	StartedAt  time.Time
	FinishedAt time.Time
	Project    Project
	TransferID string
}

func NewHistory(bytes int64, status string, startedAt, finishedAt time.Time, project Project, transferID string) History {
	history := History{
		Bytes:      bytes,
		Status:     status,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		Project:    project,
		TransferID: transferID,
	}
	return history
}

type HistoryStorage interface {
	Create(history *History) error
	Read(id string) (History, error)
	Update(history History) error
	Delete(history History) error
}
