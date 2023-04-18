package transfer

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/database"
)

type Repository interface {
	Create(transfer *Transfer) error
	Read(id string) (Transfer, error)
	List() ([]Transfer, error)
	Update(transfer Transfer) error
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

func (r *PostgresRepository) Create(transfer *Transfer) error {
	stmt := `
		INSERT INTO transfer
			(pattern, from_location_id, to_location_id, schedule_id)
		VALUES
			($1, $2, $3, $4)
		RETURNING id`

	args := []any{
		transfer.Pattern,
		transfer.FromLocationID,
		transfer.ToLocationID,
		transfer.ScheduleID,
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

func (r *PostgresRepository) Read(id string) (Transfer, error) {
	stmt := `
		SELECT
			transfer.id,
			transfer.pattern,
			transfer.from_location_id,
			transfer.to_location_id,
			transfer.schedule_id
		FROM transfer
		WHERE transfer.id = $1`

	var transfer Transfer
	dest := []any{
		&transfer.ID,
		&transfer.Pattern,
		&transfer.FromLocationID,
		&transfer.ToLocationID,
		&transfer.ScheduleID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := r.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Read(id)
		}

		return Transfer{}, err
	}

	return transfer, nil
}

func (r *PostgresRepository) List() ([]Transfer, error) {
	stmt := `
		SELECT
			transfer.id,
			transfer.pattern,
			transfer.from_location_id,
			transfer.to_location_id,
			transfer.schedule_id
		FROM transfer`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []Transfer
	for rows.Next() {
		var transfer Transfer
		dest := []any{
			&transfer.ID,
			&transfer.Pattern,
			&transfer.FromLocationID,
			&transfer.ToLocationID,
			&transfer.ScheduleID,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, database.ErrRetry) {
				return r.List()
			}

			return nil, err
		}

		transfers = append(transfers, transfer)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return transfers, nil
}

func (r *PostgresRepository) Update(transfer Transfer) error {
	stmt := `
		UPDATE transfer
		SET
			pattern = $2,
			from_location_id = $3,
			to_location_id = $4,
			schedule_id = $5
		WHERE id = $1`

	args := []any{
		transfer.ID,
		transfer.Pattern,
		transfer.FromLocationID,
		transfer.ToLocationID,
		transfer.ScheduleID,
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

func (r *PostgresRepository) Delete(transfer Transfer) error {
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
