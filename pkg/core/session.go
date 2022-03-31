package core

import (
	"time"
)

type Session struct {
	Hash    string
	Expiry  time.Time
	Account Account
}

func NewSession(hash string, expiry time.Time, account Account) Session {
	session := Session{
		Hash:    hash,
		Expiry:  expiry,
		Account: account,
	}
	return session
}
