package history

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/database"
)

type Repository interface {
	Create(history *History) error
	Read(id string) (History, error)
	List() ([]History, error)
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

func (r *PostgresRepository) Create(history *History) error {
	stmt := `
		INSERT INTO history
			(bytes, started_at, finished_at, transfer_id)
		VALUES
			($1, $2, $3, $4)
		RETURNING id`

	args := []any{
		history.Bytes,
		history.StartedAt,
		history.FinishedAt,
		history.TransferID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := r.conn.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &history.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Create(history)
		}

		return err
	}

	return nil
}

func (r *PostgresRepository) Read(id string) (History, error) {
	stmt := `
		SELECT
			history.id,
			history.bytes,
			history.started_at,
			history.finished_at,
			history.transfer_id
		FROM history
		WHERE history.id = $1`

	var history History
	dest := []any{
		&history.ID,
		&history.Bytes,
		&history.StartedAt,
		&history.FinishedAt,
		&history.TransferID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := r.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Read(id)
		}

		return History{}, err
	}

	return history, nil
}

func (r *PostgresRepository) List() ([]History, error) {
	stmt := `
		SELECT
			history.id,
			history.bytes,
			history.started_at,
			history.finished_at,
			history.transfer_id
		FROM history`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []History
	for rows.Next() {
		var history History
		dest := []any{
			&history.ID,
			&history.Bytes,
			&history.StartedAt,
			&history.FinishedAt,
			&history.TransferID,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, database.ErrRetry) {
				return r.List()
			}

			return nil, err
		}

		histories = append(histories, history)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return histories, nil
}

func (r *PostgresRepository) Delete(history History) error {
	stmt := `
		DELETE FROM history
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	err := database.Exec(r.conn, ctx, stmt, history.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Delete(history)
		}

		return err
	}

	return nil
}
