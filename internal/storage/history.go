package storage

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/postgresql"
)

type History struct {
	db postgresql.Conn
}

func NewHistory(db postgresql.Conn) *History {
	s := History{
		db: db,
	}
	return &s
}

func (s *History) Create(history *model.History) error {
	stmt := `
		INSERT INTO history
			(bytes, status, started_at, finished_at, transfer_id)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id`

	args := []any{
		history.Bytes,
		history.Status,
		history.StartedAt,
		history.FinishedAt,
		history.TransferID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, args...)
	err := postgresql.Scan(row, &history.ID)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.Create(history)
		}

		return err
	}

	return nil
}

func (s *History) Read(id string) (model.History, error) {
	stmt := `
		SELECT
			history.id,
			history.bytes,
			history.status,
			history.started_at,
			history.finished_at,
			history.transfer_id
		FROM history
		WHERE history.id = $1`

	var history model.History
	dest := []any{
		&history.ID,
		&history.Bytes,
		&history.Status,
		&history.StartedAt,
		&history.FinishedAt,
		&history.TransferID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, id)
	err := postgresql.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.Read(id)
		}

		return model.History{}, err
	}

	return history, nil
}

func (s *History) ReadAll() ([]model.History, error) {
	stmt := `
		SELECT
			history.id,
			history.bytes,
			history.status,
			history.started_at,
			history.finished_at,
			history.transfer_id
		FROM history
		LEFT JOIN transfer
			ON transfer.id = history.transfer_id`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := s.db.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []model.History
	for rows.Next() {
		var history model.History
		dest := []any{
			&history.ID,
			&history.Bytes,
			&history.Status,
			&history.StartedAt,
			&history.FinishedAt,
			&history.TransferID,
		}

		err := postgresql.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, postgresql.ErrRetry) {
				return s.ReadAll()
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

func (s *History) Delete(history model.History) error {
	stmt := `
		DELETE FROM history
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := postgresql.Exec(s.db, ctx, stmt, history.ID)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.Delete(history)
		}

		return err
	}

	return nil
}
