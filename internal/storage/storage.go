package storage

import (
	"time"

	"github.com/theandrew168/dripfile/internal/database"
)

// default query timeout
const timeout = 3 * time.Second

// aggregation of core storage types
type Storage struct {
	db database.Conn

	Project  *Project
	Account  *Account
	Session  *Session
	Location *Location
	Transfer *Transfer
	Schedule *Schedule
	History  *History
}

func New(db database.Conn) *Storage {
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
