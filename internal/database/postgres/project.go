package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/postgres"
)

type projectStorage struct {
	pool *pgxpool.Pool
}

func NewProjectStorage(pool *pgxpool.Pool) *projectStorage {
	s := projectStorage{
		pool: pool,
	}
	return &s
}

func (s *projectStorage) Create(project *core.Project) error {
	stmt := `
		INSERT INTO project
			(billing_id, billing_verified)
		VALUES
			($1, $2)
		RETURNING id`

	args := []interface{}{
		project.BillingID,
		project.BillingVerified,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.pool.QueryRow(ctx, stmt, args...)
	err := postgres.Scan(row, &project.ID)
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
			project.id,
			project.billing_id,
			project.billing_verified
		FROM project
		WHERE project.id = $1`

	var project core.Project
	dest := []interface{}{
		&project.ID,
		&project.BillingID,
		&project.BillingVerified,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.pool.QueryRow(ctx, stmt, id)
	err := postgres.Scan(row, dest...)
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
			billing_id = $2,
			billing_verified = $3
		WHERE id = $1`

	args := []interface{}{
		project.ID,
		project.BillingID,
		project.BillingVerified,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(s.pool, ctx, stmt, args...)
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

	err := postgres.Exec(s.pool, ctx, stmt, project.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Delete(project)
		}

		return err
	}

	return nil
}
