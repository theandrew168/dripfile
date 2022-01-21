package core

import (
	"time"
)

type Session struct {
	// readonly (from database, after creation)
	ID int64

	Expiry  time.Time
	Account Account
}

func NewSession(expiry time.Time, account Account) Session {
	session := Session{
		Expiry:  expiry,
		Account: account,
	}
	return session
}

type SessionStorage interface {
	Create(session *Session) error
	Read(id int64) (Session, error)
	Update(session Session) error
	Delete(session Session) error
}
