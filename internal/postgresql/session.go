package postgresql

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type sessionStorage struct {
	conn *pgxpool.Pool
}

func NewSessionStorage(conn *pgxpool.Pool) core.SessionStorage {
	s := sessionStorage{
		conn: conn,
	}
	return &s
}

func (s *sessionStorage) Create(session *core.Session) error {
	return nil
}

func (s *sessionStorage) Read(id int64) (core.Session, error) {
	return core.Session{}, nil
}

func (s *sessionStorage) Update(session core.Session) error {
	return nil
}

func (s *sessionStorage) Delete(session core.Session) error {
	return nil
}
