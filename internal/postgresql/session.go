package postgresql

import (
	"context"
	"errors"

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
	stmt := `
		INSERT INTO session
			(id, expiry, account_id)
		VALUES
			($1, $2, $3)`

	args := []interface{}{
		session.ID,
		session.Expiry,
		session.Account.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := exec(s.conn, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Create(session)
		}

		return err
	}

	return nil
}

func (s *sessionStorage) Read(id string) (core.Session, error) {
	stmt := `
		SELECT
			session.id,
			session.expiry,
			account.id,
			account.email,
			account.username,
			account.password,
			account.verified
		FROM session
		INNER JOIN account
			ON account.id = session.account_id
		WHERE session.id = $1`

	var session core.Session
	dest := []interface{}{
		&session.ID,
		&session.Expiry,
		&session.Account.ID,
		&session.Account.Email,
		&session.Account.Username,
		&session.Account.Password,
		&session.Account.Verified,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.conn.QueryRow(ctx, stmt, id)
	err := scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Read(id)
		}

		return core.Session{}, err
	}

	return session, nil
}

func (s *sessionStorage) Delete(session core.Session) error {
	stmt := `
		DELETE FROM session
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := exec(s.conn, ctx, stmt, session.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Delete(session)
		}

		return err
	}

	return nil
}
