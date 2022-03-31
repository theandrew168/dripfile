package database

import (
	"time"

	"github.com/theandrew168/dripfile/pkg/postgres"
)

// default query timeout
var queryTimeout = 3 * time.Second

// aggregation of core storage types
type Storage struct {
	db postgres.Database

	Project  *ProjectStorage
	Account  *AccountStorage
	Session  *SessionStorage
	Location *LocationStorage
	Transfer *TransferStorage
	Schedule *ScheduleStorage
	Job      *JobStorage
	History  *HistoryStorage
}

func NewStorage(db postgres.Database) *Storage {
	s := Storage{
		db: db,

		Project:  NewProjectStorage(db),
		Account:  NewAccountStorage(db),
		Session:  NewSessionStorage(db),
		Location: NewLocationStorage(db),
		Transfer: NewTransferStorage(db),
		Schedule: NewScheduleStorage(db),
		Job:      NewJobStorage(db),
		History:  NewHistoryStorage(db),
	}
	return &s
}

func (s *Storage) WithTransaction(func() error) error {
	return nil
}
