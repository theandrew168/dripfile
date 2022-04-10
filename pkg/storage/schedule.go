package storage

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/postgres"
)

type Schedule struct {
	pg postgres.Interface
}

func NewSchedule(pg postgres.Interface) *Schedule {
	s := Schedule{
		pg: pg,
	}
	return &s
}

func (s *Schedule) Create(schedule *core.Schedule) error {
	stmt := `
		INSERT INTO schedule
			(name, expr, project_id)
		VALUES
			($1, $2, $3)
		RETURNING id`

	args := []interface{}{
		schedule.Name,
		schedule.Expr,
		schedule.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.pg.QueryRow(ctx, stmt, args...)
	err := postgres.Scan(row, &schedule.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Create(schedule)
		}

		return err
	}

	return nil
}

func (s *Schedule) Read(id string) (core.Schedule, error) {
	stmt := `
		SELECT
			schedule.id,
			schedule.name,
			schedule.expr,
			project.id,
			project.customer_id,
			project.subscription_item_id
		FROM schedule
		INNER JOIN project
			ON project.id = schedule.project_id
		WHERE schedule.id = $1`

	var schedule core.Schedule
	dest := []interface{}{
		&schedule.ID,
		&schedule.Name,
		&schedule.Expr,
		&schedule.Project.ID,
		&schedule.Project.CustomerID,
		&schedule.Project.SubscriptionItemID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.pg.QueryRow(ctx, stmt, id)
	err := postgres.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Read(id)
		}

		return core.Schedule{}, err
	}

	return schedule, nil
}

func (s *Schedule) Update(schedule core.Schedule) error {
	return nil
}

func (s *Schedule) Delete(schedule core.Schedule) error {
	stmt := `
		DELETE FROM schedule
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(s.pg, ctx, stmt, schedule.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Delete(schedule)
		}

		return err
	}

	return nil
}

func (s *Schedule) ReadAllByProject(project core.Project) ([]core.Schedule, error) {
	stmt := `
		SELECT
			schedule.id,
			schedule.name,
			schedule.expr,
			project.id,
			project.customer_id,
			project.subscription_item_id
		FROM schedule
		INNER JOIN project
			ON project.id = schedule.project_id
		WHERE project.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := s.pg.Query(ctx, stmt, project.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []core.Schedule
	for rows.Next() {
		var schedule core.Schedule
		dest := []interface{}{
			&schedule.ID,
			&schedule.Name,
			&schedule.Expr,
			&schedule.Project.ID,
			&schedule.Project.CustomerID,
			&schedule.Project.SubscriptionItemID,
		}

		err := postgres.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, core.ErrRetry) {
				return s.ReadAllByProject(project)
			}

			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return schedules, nil
}
