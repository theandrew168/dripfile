package postgresql

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
			(email, username, password, verified, role)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id`

	args := []interface{}{
		account.Email,
		account.Username,
		account.Password,
		account.Verified,
		account.Role,
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

func (s *accountStorage) Read(id int64) (core.Account, error) {
	stmt := `
		SELECT
			account.id,
			account.email,
			account.username,
			account.password,
			account.verified,
			account.role
		FROM account
		WHERE account.id = $1`

	var account core.Account
	dest := []interface{}{
		&account.ID,
		&account.Email,
		&account.Username,
		&account.Password,
		&account.Verified,
		&account.Role,
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
			verified = $5,
			role = $6
		WHERE id = $1`

	args := []interface{}{
		account.ID,
		account.Email,
		account.Username,
		account.Password,
		account.Verified,
		account.Role,
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
			account.verified,
			account.role
		FROM account
		WHERE account.email = $1`

	var account core.Account
	dest := []interface{}{
		&account.ID,
		&account.Email,
		&account.Username,
		&account.Password,
		&account.Verified,
		&account.Role,
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

func (s *accountStorage) ReadManyByProject(project core.Project) ([]core.Account, error) {
	stmt := `
		SELECT
			account.id,
			account.email,
			account.username,
			account.password,
			account.verified,
			account.role
		FROM account
		INNER JOIN project
			ON project.id = account.project_id
		WHERE project.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, project.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []core.Account
	for rows.Next() {
		var account core.Account
		dest := []interface{}{
			&account.ID,
			&account.Email,
			&account.Username,
			&account.Password,
			&account.Verified,
			&account.Role,
		}

		err := scan(rows, dest...)
		if err != nil {
			if errors.Is(err, core.ErrRetry) {
				return s.ReadManyByProject(project)
			}

			return nil, err
		}

		accounts = append(accounts, account)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return accounts, nil
}
