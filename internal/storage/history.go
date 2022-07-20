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
	return model.History{}, nil
}

func (s *History) ReadAllByProject(project model.Project) ([]model.History, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := s.db.Query(ctx, stmt, project.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []model.History
	for rows.Next() {
		var history model.History
		dest := []interface{}{
			&history.ID,
			&history.Bytes,
			&history.StartedAt,
			&history.FinishedAt,
			&history.TransferID,
			&history.Project.ID,
		}

		err := postgresql.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, postgresql.ErrRetry) {
				return s.ReadAllByProject(project)
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
