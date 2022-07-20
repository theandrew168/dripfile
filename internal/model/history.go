package model

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
	TransferID string
	Project    Project
}

func NewHistory(bytes int64, status string, startedAt, finishedAt time.Time, transferID string, project Project) History {
	history := History{
		Bytes:      bytes,
		Status:     status,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		TransferID: transferID,
		Project:    project,
	}
	return history
}
