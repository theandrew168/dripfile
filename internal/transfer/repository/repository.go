package transfer

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/transfer"
)

type Repository interface {
	Create(transfer *transfer.Transfer) error
	Read(id string) (transfer.Transfer, error)
	List() ([]transfer.Transfer, error)
	Update(transfer transfer.Transfer) error
	Delete(id string) error
}

type PostgresRepository struct {
	conn database.Conn
}

func NewPostgresRepository(conn database.Conn) *PostgresRepository {
	r := PostgresRepository{
		conn: conn,
	}
	return &r
}

func (r *PostgresRepository) Create(transfer *transfer.Transfer) error {
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

	row := r.conn.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &transfer.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Create(transfer)
		}

		return err
	}

	return nil
}

func (r *PostgresRepository) Read(id string) (transfer.Transfer, error) {
	stmt := `
		SELECT
			transfer.id,
			transfer.pattern,
			transfer.from_location_id,
			transfer.to_location_id
		FROM transfer
		WHERE transfer.id = $1`

	var m transfer.Transfer
	dest := []any{
		&m.ID,
		&m.Pattern,
		&m.FromLocationID,
		&m.ToLocationID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := r.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Read(id)
		}

		return transfer.Transfer{}, err
	}

	return m, nil
}

func (r *PostgresRepository) List() ([]transfer.Transfer, error) {
	stmt := `
		SELECT
			transfer.id,
			transfer.pattern,
			transfer.from_location_id,
			transfer.to_location_id
		FROM transfer`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ms []transfer.Transfer
	for rows.Next() {
		var m transfer.Transfer
		dest := []any{
			&m.ID,
			&m.Pattern,
			&m.FromLocationID,
			&m.ToLocationID,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, database.ErrRetry) {
				return r.List()
			}

			return nil, err
		}

		ms = append(ms, m)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return ms, nil
}

func (r *PostgresRepository) Update(transfer transfer.Transfer) error {
	stmt := `
		UPDATE transfer
		SET
			pattern = $2,
			from_location_id = $3,
			to_location_id = $4
		WHERE id = $1`

	args := []any{
		transfer.ID,
		transfer.Pattern,
		transfer.FromLocationID,
		transfer.ToLocationID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	err := database.Exec(r.conn, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Update(transfer)
		}

		return err
	}

	return nil
}

func (r *PostgresRepository) Delete(transfer transfer.Transfer) error {
	stmt := `
		DELETE FROM transfer
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	err := database.Exec(r.conn, ctx, stmt, transfer.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Delete(transfer)
		}

		return err
	}

	return nil
}
