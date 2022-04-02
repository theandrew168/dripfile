package database

import (
	"time"

	"github.com/theandrew168/dripfile/pkg/postgres"
)

// default query timeout
var queryTimeout = 3 * time.Second

// aggregation of core storage types
type Storage struct {
	pg postgres.Interface

	Project  *ProjectStorage
	Account  *AccountStorage
	Session  *SessionStorage
	Location *LocationStorage
	Transfer *TransferStorage
	Schedule *ScheduleStorage
	Job      *JobStorage
	History  *HistoryStorage
}

func NewStorage(pg postgres.Interface) *Storage {
	s := Storage{
		pg: pg,

		Project:  NewProjectStorage(pg),
		Account:  NewAccountStorage(pg),
		Session:  NewSessionStorage(pg),
		Location: NewLocationStorage(pg),
		Transfer: NewTransferStorage(pg),
		Schedule: NewScheduleStorage(pg),
		Job:      NewJobStorage(pg),
		History:  NewHistoryStorage(pg),
	}
	return &s
}

func (s *Storage) WithTransaction(func() error) error {
	return nil
}
