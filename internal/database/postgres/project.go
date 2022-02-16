package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type projectStorage struct {
	conn *pgxpool.Pool
}

func NewProjectStorage(conn *pgxpool.Pool) *projectStorage {
	s := projectStorage{
		conn: conn,
	}
	return &s
}

func (s *projectStorage) Create(project *core.Project) error {
	stmt := `
		INSERT INTO project
		DEFAULT VALUES
		RETURNING id`

	args := []interface{}{}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.conn.QueryRow(ctx, stmt, args...)
	err := scan(row, &project.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Create(project)
		}

		return err
	}

	return nil
}

func (s *projectStorage) Read(id string) (core.Project, error) {
	stmt := `
		SELECT
			project.id
		FROM project
		WHERE project.id = $1`

	var project core.Project
	dest := []interface{}{
		&project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.conn.QueryRow(ctx, stmt, id)
	err := scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Read(id)
		}

		return core.Project{}, err
	}

	return project, nil
}

func (s *projectStorage) Delete(project core.Project) error {
	stmt := `
		DELETE FROM project
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := exec(s.conn, ctx, stmt, project.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Delete(project)
		}

		return err
	}

	return nil
}
