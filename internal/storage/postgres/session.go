package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type sessionStorage struct {
	conn *pgxpool.Pool
}

func NewSession(conn *pgxpool.Pool) *sessionStorage {
	s := sessionStorage{
		conn: conn,
	}
	return &s
}

func (s *sessionStorage) Create(session *core.Session) error {
	stmt := `
		INSERT INTO session
			(hash, expiry, account_id)
		VALUES
			($1, $2, $3)`

	args := []interface{}{
		session.Hash,
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

func (s *sessionStorage) Read(hash string) (core.Session, error) {
	stmt := `
		SELECT
			session.hash,
			session.expiry,
			account.id,
			account.email,
			account.username,
			account.password,
			account.role,
			account.verified,
			project.id
		FROM session
		INNER JOIN account
			ON account.id = session.account_id
		INNER JOIN project
			ON project.id = account.project_id
		WHERE session.hash = $1`

	var session core.Session
	dest := []interface{}{
		&session.Hash,
		&session.Expiry,
		&session.Account.ID,
		&session.Account.Email,
		&session.Account.Username,
		&session.Account.Password,
		&session.Account.Role,
		&session.Account.Verified,
		&session.Account.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.conn.QueryRow(ctx, stmt, hash)
	err := scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Read(hash)
		}

		return core.Session{}, err
	}

	return session, nil
}

func (s *sessionStorage) Delete(session core.Session) error {
	stmt := `
		DELETE FROM session
		WHERE hash = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := exec(s.conn, ctx, stmt, session.Hash)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Delete(session)
		}

		return err
	}

	return nil
}
