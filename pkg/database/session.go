package database

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/postgres"
)

type SessionStorage struct {
	db postgres.Database
}

func NewSessionStorage(db postgres.Database) *SessionStorage {
	s := SessionStorage{
		db: db,
	}
	return &s
}

func (s *SessionStorage) Create(session *core.Session) error {
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

	err := postgres.Exec(s.db, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Create(session)
		}

		return err
	}

	return nil
}

func (s *SessionStorage) Read(hash string) (core.Session, error) {
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
			project.id,
			project.customer_id,
			project.subscription_item_id
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
		&session.Account.Project.CustomerID,
		&session.Account.Project.SubscriptionItemID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, hash)
	err := postgres.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Read(hash)
		}

		return core.Session{}, err
	}

	return session, nil
}

func (s *SessionStorage) Delete(session core.Session) error {
	stmt := `
		DELETE FROM session
		WHERE hash = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(s.db, ctx, stmt, session.Hash)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Delete(session)
		}

		return err
	}

	return nil
}

func (s *SessionStorage) DeleteExpired() error {
	stmt := `
		DELETE FROM session
		WHERE expiry <= NOW()`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(s.db, ctx, stmt)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.DeleteExpired()
		}

		return err
	}

	return nil
}
