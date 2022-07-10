package storage

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
)

type Account struct {
	db database.Conn
}

func NewAccount(db database.Conn) *Account {
	s := Account{
		db: db,
	}
	return &s
}

func (s *Account) Create(account *core.Account) error {
	stmt := `
		INSERT INTO account
			(email, password, role, verified, project_id)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id`

	args := []interface{}{
		account.Email,
		account.Password,
		account.Role,
		account.Verified,
		account.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &account.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Create(account)
		}

		return err
	}

	return nil
}

func (s *Account) Read(id string) (core.Account, error) {
	stmt := `
		SELECT
			account.id,
			account.email,
			account.password,
			account.role,
			account.verified,
			project.id
		FROM account
		INNER JOIN project
			ON project.id = account.project_id
		WHERE account.id = $1`

	var account core.Account
	dest := []interface{}{
		&account.ID,
		&account.Email,
		&account.Password,
		&account.Role,
		&account.Verified,
		&account.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Read(id)
		}

		return core.Account{}, err
	}

	return account, nil
}

func (s *Account) Update(account core.Account) error {
	stmt := `
		UPDATE account
		SET
			email = $2,
			password = $3,
			role = $4,
			verified = $5
		WHERE id = $1`

	args := []interface{}{
		account.ID,
		account.Email,
		account.Password,
		account.Role,
		account.Verified,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := database.Exec(s.db, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Update(account)
		}

		return err
	}

	return nil
}

func (s *Account) Delete(account core.Account) error {
	stmt := `
		DELETE FROM account
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := database.Exec(s.db, ctx, stmt, account.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Delete(account)
		}

		return err
	}

	return nil
}

func (s *Account) ReadByEmail(email string) (core.Account, error) {
	stmt := `
		SELECT
			account.id,
			account.email,
			account.password,
			account.role,
			account.verified,
			project.id
		FROM account
		INNER JOIN project
			ON project.id = account.project_id
		WHERE account.email = $1`

	var account core.Account
	dest := []interface{}{
		&account.ID,
		&account.Email,
		&account.Password,
		&account.Role,
		&account.Verified,
		&account.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, email)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.ReadByEmail(email)
		}

		return core.Account{}, err
	}

	return account, nil
}

func (s *Account) CountByProject(project core.Project) (int, error) {
	stmt := `
		SELECT
			count(*)
		FROM account
		INNER JOIN project
			ON project.id = account.project_id
		WHERE project.id = $1`

	var count int

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, project.ID)
	err := database.Scan(row, &count)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.CountByProject(project)
		}

		return 0, err
	}

	return count, nil
}
