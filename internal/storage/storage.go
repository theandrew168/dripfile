package storage

import (
	"time"

	"github.com/theandrew168/dripfile/internal/postgresql"
)

// default query timeout
const timeout = 3 * time.Second

// aggregation of core storage types
type Storage struct {
	db postgresql.Conn

	Account  *Account
	Session  *Session
	Location *Location
	Transfer *Transfer
	Schedule *Schedule
	History  *History
}

func New(db postgresql.Conn) *Storage {
	s := Storage{
		db: db,

		Account:  NewAccount(db),
		Session:  NewSession(db),
		Location: NewLocation(db),
		Transfer: NewTransfer(db),
		Schedule: NewSchedule(db),
		History:  NewHistory(db),
	}
	return &s
}
