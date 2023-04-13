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

	Account  *Account
	Session  *Session
	Location *Location
}

func New(db database.Conn) *Storage {
	s := Storage{
		db: db,

		Account:  NewAccount(db),
		Session:  NewSession(db),
		Location: NewLocation(db),
	}
	return &s
}
