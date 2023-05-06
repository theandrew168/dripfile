package history

import "time"

type History struct {
	// readonly (from database, after creation)
	ID string

	Bytes      int64
	StartedAt  time.Time
	FinishedAt time.Time
	TransferID string
}

func New(bytes int64, startedAt, finishedAt time.Time, transferID string) History {
	history := History{
		Bytes:      bytes,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		TransferID: transferID,
	}
	return history
}

type Repository interface {
	Create(history *History) error
	Read(id string) (History, error)
	List() ([]History, error)
	Delete(id string) error
}
