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
}

func NewHistory(bytes int64, status string, startedAt, finishedAt time.Time, transferID string) History {
	history := History{
		Bytes:      bytes,
		Status:     status,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		TransferID: transferID,
	}
	return history
}
