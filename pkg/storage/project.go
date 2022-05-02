package storage

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/database"
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
			(customer_id, subscription_item_id)
		VALUES
			($1, $2)
		RETURNING id`

	args := []interface{}{
		project.CustomerID,
		project.SubscriptionItemID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &project.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Create(project)
		}

		return err
	}

	return nil
}

func (s *Project) Read(id string) (core.Project, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Read(id)
		}

		return core.Project{}, err
	}

	return project, nil
}

func (s *Project) Update(project core.Project) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := database.Exec(s.db, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Update(project)
		}

		return err
	}

	return nil
}

func (s *Project) Delete(project core.Project) error {
	stmt := `
		DELETE FROM project
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := database.Exec(s.db, ctx, stmt, project.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Delete(project)
		}

		return err
	}

	return nil
}

func (s *Project) ReadAll() ([]core.Project, error) {
	stmt := `
		SELECT
			project.id,
			project.customer_id,
			project.subscription_item_id
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
			&project.CustomerID,
			&project.SubscriptionItemID,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, core.ErrRetry) {
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
