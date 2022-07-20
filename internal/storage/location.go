package storage

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/postgresql"
)

type Location struct {
	db postgresql.Conn
}

func NewLocation(db postgresql.Conn) *Location {
	s := Location{
		db: db,
	}
	return &s
}

func (s *Location) Create(location *model.Location) error {
	stmt := `
		INSERT INTO location
			(kind, name, info, project_id)
		VALUES
			($1, $2, $3, $4)
		RETURNING id`

	args := []interface{}{
		location.Kind,
		location.Name,
		location.Info,
		location.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, args...)
	err := postgresql.Scan(row, &location.ID)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.Create(location)
		}

		return err
	}

	return nil
}

func (s *Location) Read(id string) (model.Location, error) {
	stmt := `
		SELECT
			location.id,
			location.kind,
			location.name,
			location.info,
			project.id
		FROM location
		INNER JOIN project
			ON project.id = location.project_id
		WHERE location.id = $1`

	var location model.Location
	dest := []interface{}{
		&location.ID,
		&location.Kind,
		&location.Name,
		&location.Info,
		&location.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, id)
	err := postgresql.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.Read(id)
		}

		return model.Location{}, err
	}

	return location, nil
}

func (s *Location) Update(location model.Location) error {
	stmt := `
		UPDATE location
		SET
			kind = $2,
			name = $3,
			info = $4
		WHERE id = $1`

	args := []interface{}{
		location.ID,
		location.Kind,
		location.Name,
		location.Info,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := postgresql.Exec(s.db, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.Update(location)
		}

		return err
	}

	return nil
}

func (s *Location) Delete(location model.Location) error {
	stmt := `
		DELETE FROM location
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := postgresql.Exec(s.db, ctx, stmt, location.ID)
	if err != nil {
		if errors.Is(err, postgresql.ErrRetry) {
			return s.Delete(location)
		}

		return err
	}

	return nil
}

func (s *Location) ReadAllByProject(project model.Project) ([]model.Location, error) {
	stmt := `
		SELECT
			location.id,
			location.kind,
			location.name,
			location.info,
			project.id
		FROM location
		INNER JOIN project
			ON project.id = location.project_id
		WHERE project.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := s.db.Query(ctx, stmt, project.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []model.Location
	for rows.Next() {
		var location model.Location
		dest := []interface{}{
			&location.ID,
			&location.Kind,
			&location.Name,
			&location.Info,
			&location.Project.ID,
		}

		err := postgresql.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, postgresql.ErrRetry) {
				return s.ReadAllByProject(project)
			}

			return nil, err
		}

		locations = append(locations, location)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return locations, nil
}
