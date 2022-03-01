package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/postgres"
)

type locationStorage struct {
	pool *pgxpool.Pool
}

func NewLocationStorage(pool *pgxpool.Pool) *locationStorage {
	s := locationStorage{
		pool: pool,
	}
	return &s
}

// TODO: return zero, one, or many results (only last op?)
type Operation struct {
	stmt string
	args []interface{}
	dest []interface{}
}

/*

// return zero
create := Operation{
	stmt: "INSERT INTO location ... ($1, $2, $3) ...",
	args: []interface{
		location.Kind,
		location.Info,
		location.Project.ID,
	},
	dest: []interface{
		&location.ID,
	},
}

// return one
read := Operation{
	stmt: "SELECT FROM location ... WHERE location.id = $1",
	args: []interface{
		location.ID,
	},
	dest: []interface{
		&location.ID,
		&location.Kind,
		&location.Info,
		&location.Project.ID,
	},
}

*/

func (s *locationStorage) Create(location *core.Location) error {
	stmt := `
		INSERT INTO location
			(kind, info, project_id)
		VALUES
			($1, $2, $3)
		RETURNING id`

	args := []interface{}{
		location.Kind,
		location.Info,
		location.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.pool.QueryRow(ctx, stmt, args...)
	err := postgres.Scan(row, &location.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Create(location)
		}

		return err
	}

	return nil
}

func (s *locationStorage) Read(id string) (core.Location, error) {
	stmt := `
		SELECT
			location.id,
			location.kind,
			location.info,
			project.id
		FROM location
		INNER JOIN project
			ON project.id = location.project_id
		WHERE location.id = $1`

	var location core.Location
	dest := []interface{}{
		&location.ID,
		&location.Kind,
		&location.Info,
		&location.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.pool.QueryRow(ctx, stmt, id)
	err := postgres.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Read(id)
		}

		return core.Location{}, err
	}

	return location, nil
}

func (s *locationStorage) Update(location core.Location) error {
	stmt := `
		UPDATE location
		SET
			kind = $2,
			info = $3
		WHERE id = $1`

	args := []interface{}{
		location.ID,
		location.Kind,
		location.Info,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(s.pool, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Update(location)
		}

		return err
	}

	return nil
}

func (s *locationStorage) Delete(location core.Location) error {
	stmt := `
		DELETE FROM location
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(s.pool, ctx, stmt, location.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Delete(location)
		}

		return err
	}

	return nil
}

func (s *locationStorage) ReadManyByProject(project core.Project) ([]core.Location, error) {
	stmt := `
		SELECT
			location.id,
			location.kind,
			location.info,
			project.id
		FROM location
		INNER JOIN project
			ON project.id = location.project_id
		WHERE project.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := s.pool.Query(ctx, stmt, project.ID)
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
			&location.Info,
			&location.Project.ID,
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
