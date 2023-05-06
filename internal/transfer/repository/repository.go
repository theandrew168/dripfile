package transfer

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/transfer"
)

type Repository struct {
	conn database.Conn
}

func New(conn database.Conn) *Repository {
	repo := Repository{
		conn: conn,
	}
	return &repo
}

func (repo *Repository) Create(transfer *transfer.Transfer) error {
	stmt := `
		INSERT INTO transfer
			(pattern, from_location_id, to_location_id)
		VALUES
			($1, $2, $3)
		RETURNING id`

	args := []any{
		transfer.Pattern,
		transfer.FromLocationID,
		transfer.ToLocationID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := repo.conn.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &transfer.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return repo.Create(transfer)
		}

		return err
	}

	return nil
}

func (repo *Repository) Read(id string) (transfer.Transfer, error) {
	stmt := `
		SELECT
			id,
			pattern,
			from_location_id,
			to_location_id
		FROM transfer
		WHERE id = $1`

	var t transfer.Transfer
	dest := []any{
		&t.ID,
		&t.Pattern,
		&t.FromLocationID,
		&t.ToLocationID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := repo.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return repo.Read(id)
		}

		return transfer.Transfer{}, err
	}

	return t, nil
}

func (repo *Repository) List() ([]transfer.Transfer, error) {
	stmt := `
		SELECT
			id,
			pattern,
			from_location_id,
			to_location_id
		FROM transfer`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ts []transfer.Transfer
	for rows.Next() {
		var t transfer.Transfer
		dest := []any{
			&t.ID,
			&t.Pattern,
			&t.FromLocationID,
			&t.ToLocationID,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, database.ErrRetry) {
				return repo.List()
			}

			return nil, err
		}

		ts = append(ts, t)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return ts, nil
}

// NOTE: differs from List only in the stmt
func (repo *Repository) ListByLocationID(locationID string) ([]transfer.Transfer, error) {
	stmt := `
		SELECT
			id,
			pattern,
			from_location_id,
			to_location_id
		FROM transfer
		WHERE from_location_id = $1
		   OR to_location_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ts []transfer.Transfer
	for rows.Next() {
		var t transfer.Transfer
		dest := []any{
			&t.ID,
			&t.Pattern,
			&t.FromLocationID,
			&t.ToLocationID,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, database.ErrRetry) {
				return repo.List()
			}

			return nil, err
		}

		ts = append(ts, t)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return ts, nil
}

func (repo *Repository) Update(transfer transfer.Transfer) error {
	stmt := `
		UPDATE transfer
		SET
			pattern = $2,
			from_location_id = $3,
			to_location_id = $4
		WHERE id = $1
		RETURNING id`

	args := []any{
		transfer.ID,
		transfer.Pattern,
		transfer.FromLocationID,
		transfer.ToLocationID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	var updatedID string
	row := repo.conn.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &updatedID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return repo.Update(transfer)
		}

		return err
	}

	return nil
}

func (repo *Repository) Delete(id string) error {
	stmt := `
		DELETE FROM transfer
		WHERE id = $1
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	var deletedID string
	row := repo.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, &deletedID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return repo.Delete(id)
		}

		return err
	}

	return nil
}
