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

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	var account core.Account
	row := s.conn.QueryRow(ctx, stmt, id)
	err := scan(
		row,
		&account.ID,
		&account.Email,
		&account.Username,
		&account.Password,
		&account.Verified,
		&account.Role,
	)
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
			return s.Update(account)
		}

		return err
	}

	return nil
}
