package storage

import (
	"time"

	"github.com/theandrew168/dripfile/pkg/database"
)

// default query timeout
var queryTimeout = 3 * time.Second

// aggregation of core storage types
type Storage struct {
	db database.Interface

	Project  *Project
	Account  *Account
	Session  *Session
	Location *Location
	Transfer *Transfer
	Schedule *Schedule
	History  *History
}

func New(db database.Interface) *Storage {
	s := Storage{
		db: db,

		Project:  NewProject(db),
		Account:  NewAccount(db),
		Session:  NewSession(db),
		Location: NewLocation(db),
		Transfer: NewTransfer(db),
		Schedule: NewSchedule(db),
		History:  NewHistory(db),
	}
	return &s
}

func (s *Storage) WithTransaction(func() error) error {
	return nil
}
