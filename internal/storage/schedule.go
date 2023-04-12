package storage

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/model"
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

func (s *Schedule) Create(schedule *model.Schedule) error {
	stmt := `
		INSERT INTO schedule
			(name, expr)
		VALUES
			($1, $2)
		RETURNING id`

	args := []any{
		schedule.Name,
		schedule.Expr,
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

func (s *Schedule) Read(id string) (model.Schedule, error) {
	stmt := `
		SELECT
			schedule.id,
			schedule.name,
			schedule.expr
		FROM schedule
		WHERE schedule.id = $1`

	var schedule model.Schedule
	dest := []any{
		&schedule.ID,
		&schedule.Name,
		&schedule.Expr,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Read(id)
		}

		return model.Schedule{}, err
	}

	return schedule, nil
}

func (s *Schedule) Update(schedule model.Schedule) error {
	return nil
}

func (s *Schedule) Delete(schedule model.Schedule) error {
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

func (s *Schedule) ReadAll() ([]model.Schedule, error) {
	stmt := `
		SELECT
			schedule.id,
			schedule.name,
			schedule.expr
		FROM schedule`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := s.db.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []model.Schedule
	for rows.Next() {
		var schedule model.Schedule
		dest := []any{
			&schedule.ID,
			&schedule.Name,
			&schedule.Expr,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, database.ErrRetry) {
				return s.ReadAll()
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
