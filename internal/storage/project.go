package storage

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
)

type Project struct {
	db database.Conn
}

func NewProject(db database.Conn) *Project {
	s := Project{
		db: db,
	}
	return &s
}

func (s *Project) Create(project *core.Project) error {
	stmt := `
		INSERT INTO project
		VALUES
			(default)
		RETURNING id`

	args := []interface{}{}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &project.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Create(project)
		}

		return err
	}

	return nil
}

func (s *Project) Read(id string) (core.Project, error) {
	stmt := `
		SELECT
			project.id
		FROM project
		WHERE project.id = $1`

	var project core.Project
	dest := []interface{}{
		&project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Read(id)
		}

		return core.Project{}, err
	}

	return project, nil
}

func (s *Project) Delete(project core.Project) error {
	stmt := `
		DELETE FROM project
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := database.Exec(s.db, ctx, stmt, project.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Delete(project)
		}

		return err
	}

	return nil
}

func (s *Project) ReadAll() ([]core.Project, error) {
	stmt := `
		SELECT
			project.id
		FROM project`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := s.db.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []core.Project
	for rows.Next() {
		var project core.Project
		dest := []interface{}{
			&project.ID,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, database.ErrRetry) {
				return s.ReadAll()
			}

			return nil, err
		}

		projects = append(projects, project)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return projects, nil
}
