package database

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/postgres"
)

type LocationStorage struct {
	db postgres.Database
}

func NewLocationStorage(db postgres.Database) *LocationStorage {
	s := LocationStorage{
		db: db,
	}
	return &s
}

func (s *LocationStorage) Create(location *core.Location) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, args...)
	err := postgres.Scan(row, &location.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Create(location)
		}

		return err
	}

	return nil
}

func (s *LocationStorage) Read(id string) (core.Location, error) {
	stmt := `
		SELECT
			location.id,
			location.kind,
			location.name,
			location.info,
			project.id,
			project.billing_id,
			project.billing_verified
		FROM location
		INNER JOIN project
			ON project.id = location.project_id
		WHERE location.id = $1`

	var location core.Location
	dest := []interface{}{
		&location.ID,
		&location.Kind,
		&location.Name,
		&location.Info,
		&location.Project.ID,
		&location.Project.BillingID,
		&location.Project.BillingVerified,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, id)
	err := postgres.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Read(id)
		}

		return core.Location{}, err
	}

	return location, nil
}

func (s *LocationStorage) Update(location core.Location) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(s.db, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Update(location)
		}

		return err
	}

	return nil
}

func (s *LocationStorage) Delete(location core.Location) error {
	stmt := `
		DELETE FROM location
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(s.db, ctx, stmt, location.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Delete(location)
		}

		return err
	}

	return nil
}

func (s *LocationStorage) ReadManyByProject(project core.Project) ([]core.Location, error) {
	stmt := `
		SELECT
			location.id,
			location.kind,
			location.name,
			location.info,
			project.id,
			project.billing_id,
			project.billing_verified
		FROM location
		INNER JOIN project
			ON project.id = location.project_id
		WHERE project.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := s.db.Query(ctx, stmt, project.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []core.Location
	for rows.Next() {
		var location core.Location
		dest := []interface{}{
			&location.ID,
			&location.Kind,
			&location.Name,
			&location.Info,
			&location.Project.ID,
			&location.Project.BillingID,
			&location.Project.BillingVerified,
		}

		err := postgres.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, core.ErrRetry) {
				return s.ReadManyByProject(project)
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
