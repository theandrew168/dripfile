package storage

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
)

type Session struct {
	db database.Conn
}

func NewSession(db database.Conn) *Session {
	s := Session{
		db: db,
	}
	return &s
}

func (s *Session) Create(session *core.Session) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := database.Exec(s.db, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Create(session)
		}

		return err
	}

	return nil
}

func (s *Session) Read(hash string) (core.Session, error) {
	stmt := `
		SELECT
			session.hash,
			session.expiry,
			account.id,
			account.email,
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
		&session.Account.Password,
		&session.Account.Role,
		&session.Account.Verified,
		&session.Account.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, hash)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Read(hash)
		}

		return core.Session{}, err
	}

	return session, nil
}

func (s *Session) Delete(session core.Session) error {
	stmt := `
		DELETE FROM session
		WHERE hash = $1`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := database.Exec(s.db, ctx, stmt, session.Hash)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Delete(session)
		}

		return err
	}

	return nil
}

func (s *Session) DeleteExpired() error {
	stmt := `
		DELETE FROM session
		WHERE expiry <= NOW()`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := database.Exec(s.db, ctx, stmt)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.DeleteExpired()
		}

		return err
	}

	return nil
}
