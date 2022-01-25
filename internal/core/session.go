package core

import (
	"time"
)

type Session struct {
	ID      string
	Expiry  time.Time
	Account Account
}

func NewSession(id string, expiry time.Time, account Account) Session {
	session := Session{
		ID:      id,
		Expiry:  expiry,
		Account: account,
	}
	return session
}

type SessionStorage interface {
	Create(session *Session) error
	Read(id string) (Session, error)
	Delete(session Session) error
}
