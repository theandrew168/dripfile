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

func (s *historyStorage) ReadManyByProject(project core.Project) ([]core.History, error) {
	stmt := `
		SELECT
			history.id,
			history.bytes,
			history.started_at,
			history.finished_at,
			history.transfer_id,
			project.id
		FROM history
		INNER JOIN project
			ON project.id = history.project_id
		LEFT JOIN transfer
			ON transfer.id = history.transfer_id
		WHERE project.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := s.pool.Query(ctx, stmt, project.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []core.History
	for rows.Next() {
		var history core.History
		dest := []interface{}{
			&history.ID,
			&history.Bytes,
			&history.StartedAt,
			&history.FinishedAt,
			&history.TransferID,
			&history.Project.ID,
		}

		err := postgres.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, core.ErrRetry) {
				return s.ReadManyByProject(project)
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
