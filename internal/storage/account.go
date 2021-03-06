package storage

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/postgresql"
)

type Account struct {
	db postgresql.Conn
}

func NewAccount(db postgresql.Conn) *Account {
	s := Account{
		db: db,
	}
	return &s
}

func (s *Account) Create(account *model.Account) error {
	stmt := `
		INSERT INTO account
			(email, password, role, verified, project_id)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id`

	args := []any{
		account.Email,
		account.Password,
		account.Role,
		account.Verified,
		account.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, args...)
	err := postgresql.Scan(row, &account.ID)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.Create(account)
		}

		return err
	}

	return nil
}

func (s *Account) Read(id string) (model.Account, error) {
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

	var account model.Account
	dest := []any{
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
	err := postgresql.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.Read(id)
		}

		return model.Account{}, err
	}

	return account, nil
}

func (s *Account) Update(account model.Account) error {
	stmt := `
		UPDATE account
		SET
			email = $2,
			password = $3,
			role = $4,
			verified = $5
		WHERE id = $1`

	args := []any{
		account.ID,
		account.Email,
		account.Password,
		account.Role,
		account.Verified,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := postgresql.Exec(s.db, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.Update(account)
		}

		return err
	}

	return nil
}

func (s *Account) Delete(account model.Account) error {
	stmt := `
		DELETE FROM account
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := postgresql.Exec(s.db, ctx, stmt, account.ID)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.Delete(account)
		}

		return err
	}

	return nil
}

func (s *Account) ReadByEmail(email string) (model.Account, error) {
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

	var account model.Account
	dest := []any{
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
	err := postgresql.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.ReadByEmail(email)
		}

		return model.Account{}, err
	}

	return account, nil
}

func (s *Account) CountByProject(project model.Project) (int, error) {
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
	err := postgresql.Scan(row, &count)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.CountByProject(project)
		}

		return 0, err
	}

	return count, nil
}
