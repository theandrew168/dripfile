package storage

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/postgresql"
)

type Project struct {
	db postgresql.Conn
}

func NewProject(db postgresql.Conn) *Project {
	s := Project{
		db: db,
	}
	return &s
}

func (s *Project) Create(project *model.Project) error {
	stmt := `
		INSERT INTO project
		VALUES
			(default)
		RETURNING id`

	args := []interface{}{}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, args...)
	err := postgresql.Scan(row, &project.ID)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.Create(project)
		}

		return err
	}

	return nil
}

func (s *Project) Read(id string) (model.Project, error) {
	stmt := `
		SELECT
			project.id
		FROM project
		WHERE project.id = $1`

	var project model.Project
	dest := []interface{}{
		&project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, id)
	err := postgresql.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.Read(id)
		}

		return model.Project{}, err
	}

	return project, nil
}

func (s *Project) Delete(project model.Project) error {
	stmt := `
		DELETE FROM project
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := postgresql.Exec(s.db, ctx, stmt, project.ID)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.Delete(project)
		}

		return err
	}

	return nil
}

func (s *Project) ReadAll() ([]model.Project, error) {
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

	var projects []model.Project
	for rows.Next() {
		var project model.Project
		dest := []interface{}{
			&project.ID,
		}

		err := postgresql.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, postgresql.ErrRetry) {
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
