package storage

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/model"
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

func (s *Account) Create(account *model.Account) error {
	stmt := `
		INSERT INTO account
			(email, password, role, verified)
		VALUES
			($1, $2, $3, $4)
		RETURNING id`

	args := []any{
		account.Email,
		account.Password,
		account.Role,
		account.Verified,
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

func (s *Account) Read(id string) (model.Account, error) {
	stmt := `
		SELECT
			account.id,
			account.email,
			account.password,
			account.role,
			account.verified
		FROM account
		WHERE account.id = $1`

	var account model.Account
	dest := []any{
		&account.ID,
		&account.Email,
		&account.Password,
		&account.Role,
		&account.Verified,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
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

	err := database.Exec(s.db, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
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

	err := database.Exec(s.db, ctx, stmt, account.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
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
			account.verified
		FROM account
		WHERE account.email = $1`

	var account model.Account
	dest := []any{
		&account.ID,
		&account.Email,
		&account.Password,
		&account.Role,
		&account.Verified,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, email)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.ReadByEmail(email)
		}

		return model.Account{}, err
	}

	return account, nil
}
