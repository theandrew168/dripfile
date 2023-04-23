package repository

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/history"
)

type Repository interface {
	Create(history *history.History) error
	Read(id string) (history.History, error)
	List() ([]history.History, error)
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

func (r *PostgresRepository) Create(history *history.History) error {
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

func (r *PostgresRepository) Read(id string) (history.History, error) {
	stmt := `
		SELECT
			history.id,
			history.bytes,
			history.started_at,
			history.finished_at,
			history.transfer_id
		FROM history
		WHERE history.id = $1`

	var m history.History
	dest := []any{
		&m.ID,
		&m.Bytes,
		&m.StartedAt,
		&m.FinishedAt,
		&m.TransferID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := r.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Read(id)
		}

		return history.History{}, err
	}

	return m, nil
}

func (r *PostgresRepository) List() ([]history.History, error) {
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

	var ms []history.History
	for rows.Next() {
		var m history.History
		dest := []any{
			&m.ID,
			&m.Bytes,
			&m.StartedAt,
			&m.FinishedAt,
			&m.TransferID,
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

func (r *PostgresRepository) Delete(id string) error {
	stmt := `
		DELETE FROM history
		WHERE id = $1
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	var deletedID string
	row := r.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, &deletedID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Delete(id)
		}

		return err
	}

	return nil
}
