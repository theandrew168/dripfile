package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type Project struct {
	// readonly (from database, after creation)
	ID int64

	Name string
}

type projectStorage struct {
	conn *pgxpool.Pool
}

func NewProjectStorage(conn *pgxpool.Pool) core.ProjectStorage {
	s := projectStorage{
		conn: conn,
	}
	return &s
}

func (s *projectStorage) Create(project *core.Project) error {
	stmt := `
		INSERT INTO project
			(name)
		VALUES
			($1)
		RETURNING id`

	args := []interface{}{
		project.Name,
	}

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

func (s *projectStorage) Read(id int64) (core.Project, error) {
	stmt := `
		SELECT
			project.id,
			project.name
		FROM project
		WHERE project.id = $1`

	var project core.Project
	dest := []interface{}{
		&project.ID,
		&project.Name,
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

func (s *projectStorage) Update(project core.Project) error {
	stmt := `
		UPDATE project
		SET
			name = $2
		WHERE id = $1`

	args := []interface{}{
		project.ID,
		project.Name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := exec(s.conn, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Update(project)
		}

		return err
	}

	return nil
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
