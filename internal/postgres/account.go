package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type accountStorage struct {
	conn *pgxpool.Pool
}

func NewAccountStorage(conn *pgxpool.Pool) core.AccountStorage {
	s := accountStorage{
		conn: conn,
	}
	return &s
}

func (s *accountStorage) Create(account *core.Account) error {
	stmt := `
		INSERT INTO account
			(email, username, password, role, verified, project_id)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING id`

	args := []interface{}{
		account.Email,
		account.Username,
		account.Password,
		account.Role,
		account.Verified,
		account.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.conn.QueryRow(ctx, stmt, args...)
	err := scan(row, &account.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Create(account)
		}

		return err
	}

	return nil
}

func (s *accountStorage) Read(id string) (core.Account, error) {
	stmt := `
		SELECT
			account.id,
			account.email,
			account.username,
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
		&account.Username,
		&account.Password,
		&account.Role,
		&account.Verified,
		&account.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.conn.QueryRow(ctx, stmt, id)
	err := scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Read(id)
		}

		return core.Account{}, err
	}

	return account, nil
}

func (s *accountStorage) Update(account core.Account) error {
	stmt := `
		UPDATE account
		SET
			email = $2,
			username = $3,
			password = $4,
			role = $5,
			verified = $6
		WHERE id = $1`

	args := []interface{}{
		account.ID,
		account.Email,
		account.Username,
		account.Password,
		account.Role,
		account.Verified,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := exec(s.conn, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Update(account)
		}

		return err
	}

	return nil
}

func (s *accountStorage) Delete(account core.Account) error {
	stmt := `
		DELETE FROM account
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := exec(s.conn, ctx, stmt, account.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Delete(account)
		}

		return err
	}

	return nil
}

func (s *accountStorage) ReadByEmail(email string) (core.Account, error) {
	stmt := `
		SELECT
			account.id,
			account.email,
			account.username,
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
		&account.Username,
		&account.Password,
		&account.Role,
		&account.Verified,
		&account.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.conn.QueryRow(ctx, stmt, email)
	err := scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.ReadByEmail(email)
		}

		return core.Account{}, err
	}

	return account, nil
}

func (s *accountStorage) CountByProject(project core.Project) (int, error) {
	stmt := `
		SELECT
			count(*)
		FROM account
		INNER JOIN project
			ON project.id = account.project_id
		WHERE project.id = $1`

	var count int 

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.conn.QueryRow(ctx, stmt, project.ID)
	err := scan(row, &count)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.CountByProject(project)
		}

		return 0, err
	}

	return count, nil
}
