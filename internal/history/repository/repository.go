package repository

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/history"
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

func (repo *Repository) Create(history *history.History) error {
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

	row := repo.conn.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &history.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return repo.Create(history)
		}

		return err
	}

	return nil
}

func (repo *Repository) Read(id string) (history.History, error) {
	stmt := `
		SELECT
			id,
			bytes,
			started_at,
			finished_at,
			transfer_id
		FROM history
		WHERE id = $1`

	var h history.History
	dest := []any{
		&h.ID,
		&h.Bytes,
		&h.StartedAt,
		&h.FinishedAt,
		&h.TransferID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := repo.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return repo.Read(id)
		}

		return history.History{}, err
	}

	return h, nil
}

func (repo *Repository) List() ([]history.History, error) {
	stmt := `
		SELECT
			id,
			bytes,
			started_at,
			finished_at,
			transfer_id
		FROM history`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hs []history.History
	for rows.Next() {
		var h history.History
		dest := []any{
			&h.ID,
			&h.Bytes,
			&h.StartedAt,
			&h.FinishedAt,
			&h.TransferID,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, database.ErrRetry) {
				return repo.List()
			}

			return nil, err
		}

		hs = append(hs, h)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return hs, nil
}

func (repo *Repository) Delete(id string) error {
	stmt := `
		DELETE FROM history
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
