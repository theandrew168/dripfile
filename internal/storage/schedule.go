package storage

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
)

type Schedule struct {
	db database.Conn
}

func NewSchedule(db database.Conn) *Schedule {
	s := Schedule{
		db: db,
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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &schedule.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
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
			project.id
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
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := database.Exec(s.db, ctx, stmt, schedule.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
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
			project.id
		FROM schedule
		INNER JOIN project
			ON project.id = schedule.project_id
		WHERE project.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := s.db.Query(ctx, stmt, project.ID)
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
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, database.ErrRetry) {
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
