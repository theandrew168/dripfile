package database

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/postgres"
)

type ProjectStorage struct {
	pg postgres.Interface
}

func NewProjectStorage(pg postgres.Interface) *ProjectStorage {
	s := ProjectStorage{
		pg: pg,
	}
	return &s
}

func (s *ProjectStorage) Create(project *core.Project) error {
	stmt := `
		INSERT INTO project
			(customer_id, subscription_item_id)
		VALUES
			($1, $2)
		RETURNING id`

	args := []interface{}{
		project.CustomerID,
		project.SubscriptionItemID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.pg.QueryRow(ctx, stmt, args...)
	err := postgres.Scan(row, &project.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Create(project)
		}

		return err
	}

	return nil
}

func (s *ProjectStorage) Read(id string) (core.Project, error) {
	stmt := `
		SELECT
			project.id,
			project.customer_id,
			project.subscription_item_id
		FROM project
		WHERE project.id = $1`

	var project core.Project
	dest := []interface{}{
		&project.ID,
		&project.CustomerID,
		&project.SubscriptionItemID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.pg.QueryRow(ctx, stmt, id)
	err := postgres.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Read(id)
		}

		return core.Project{}, err
	}

	return project, nil
}

func (s *ProjectStorage) Update(project core.Project) error {
	stmt := `
		UPDATE project
		SET
			customer_id = $2,
			subscription_item_id = $3
		WHERE id = $1`

	args := []interface{}{
		project.ID,
		project.CustomerID,
		project.SubscriptionItemID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(s.pg, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Update(project)
		}

		return err
	}

	return nil
}

func (s *ProjectStorage) Delete(project core.Project) error {
	stmt := `
		DELETE FROM project
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(s.pg, ctx, stmt, project.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Delete(project)
		}

		return err
	}

	return nil
}
