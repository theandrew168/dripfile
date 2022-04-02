package storage

import (
	"time"

	"github.com/theandrew168/dripfile/pkg/postgres"
)

// default query timeout
var queryTimeout = 3 * time.Second

// aggregation of core storage types
type Storage struct {
	pg postgres.Interface

	Project  *Project
	Account  *Account
	Session  *Session
	Location *Location
	Transfer *Transfer
	Schedule *Schedule
	Job      *Job
	History  *History
}

func New(pg postgres.Interface) *Storage {
	s := Storage{
		pg: pg,

		Project:  NewProject(pg),
		Account:  NewAccount(pg),
		Session:  NewSession(pg),
		Location: NewLocation(pg),
		Transfer: NewTransfer(pg),
		Schedule: NewSchedule(pg),
		Job:      NewJob(pg),
		History:  NewHistory(pg),
	}
	return &s
}

func (s *Storage) WithTransaction(func() error) error {
	return nil
}
