package history

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidUUID = errors.New("history: invalid UUID")
)

type History struct {
	id string

	totalBytes int64
	startedAt  time.Time
	finishedAt time.Time
	transferID string

	// internal fields used for storage conflict resolution
	createdAt time.Time
	version   int
}

func New(id string, totalBytes int64, startedAt, finishedAt time.Time, transferID string) (*History, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidUUID
	}

	_, err = uuid.Parse(transferID)
	if err != nil {
		return nil, ErrInvalidUUID
	}

	h := History{
		id: id,

		totalBytes: totalBytes,
		startedAt:  startedAt,
		finishedAt: finishedAt,
		transferID: transferID,
	}
	return &h, nil
}

func (h *History) ID() string {
	return h.id
}

func (h *History) TotalBytes() int64 {
	return h.totalBytes
}

func (h *History) StartedAt() time.Time {
	return h.startedAt
}

func (h *History) FinishedAt() time.Time {
	return h.finishedAt
}

func (h *History) TransferID() string {
	return h.transferID
}
