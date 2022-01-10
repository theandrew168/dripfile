package core

import (
	"context"
	"time"
)

type History struct {
	JobID      int
	Bytes      int
	Status     string
	StartedAt  time.Time
	FinishedAt time.Time
	Project    Project

	// readonly (from database, after creation)
	ID int
}

func NewHistory(jobID, bytes int, status string, startedAt, finishedAt time.Time, project Project) History {
	history := History{
		JobID:      jobID,
		Bytes:      bytes,
		Status:     status,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		Project:    project,
	}
	return history
}

type HistoryStorage interface {
	CreateHistory(ctx context.Context, history *History) error
	ReadHistory(ctx context.Context, id int) (History, error)
	UpdateHistory(ctx context.Context, history History) error
	DeleteHistory(ctx context.Context, history History) error
}
