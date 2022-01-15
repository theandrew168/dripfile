package core

import (
	"time"
)

type History struct {
	// readonly (from database, after creation)
	ID int64

	TransferID int64
	Bytes      int64
	Status     string
	StartedAt  time.Time
	FinishedAt time.Time
	Project    Project
}

func NewHistory(transferID, bytes int64, status string, startedAt, finishedAt time.Time, project Project) History {
	history := History{
		TransferID: transferID,
		Bytes:      bytes,
		Status:     status,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		Project:    project,
	}
	return history
}

type HistoryStorage interface {
	Create(history *History) error
	Read(id int64) (History, error)
	Update(history History) error
	Delete(history History) error
}
