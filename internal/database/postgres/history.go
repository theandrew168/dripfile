package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/postgres"
)

type historyStorage struct {
	pool *pgxpool.Pool
}

func NewHistoryStorage(pool *pgxpool.Pool) *historyStorage {
	s := historyStorage{
		pool: pool,
	}
	return &s
}

func (s *historyStorage) Create(history *core.History) error {
	stmt := `
		INSERT INTO history
			(bytes, status, started_at, finished_at, transfer_id, project_id)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING id`

	args := []interface{}{
		history.Bytes,
		history.Status,
		history.StartedAt,
		history.FinishedAt,
		history.TransferID,
		history.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.pool.QueryRow(ctx, stmt, args...)
	err := postgres.Scan(row, &history.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Create(history)
		}

		return err
	}

	return nil
}

func (s *historyStorage) Read(id string) (core.History, error) {
	return core.History{}, nil
}

func (s *historyStorage) Update(history core.History) error {
	return nil
}

func (s *historyStorage) Delete(history core.History) error {
	return nil
}
