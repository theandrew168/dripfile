package schedule

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/common/database"
)

type Repository interface {
	Create(schedule *Schedule) error
	Read(id string) (Schedule, error)
	List() ([]Schedule, error)
	Update(location Schedule) error
	Delete(id string) error
}

type PostgresRepository struct {
	conn database.Conn
}

func NewPostgresRepository(conn database.Conn) *PostgresRepository {
	r := PostgresRepository{
		conn: conn,
	}
	return &r
}

func (r *PostgresRepository) Create(schedule *Schedule) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := r.conn.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &schedule.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Create(schedule)
		}

		return err
	}

	return nil
}

func (r *PostgresRepository) Read(id string) (Schedule, error) {
	stmt := `
		SELECT
			schedule.id,
			schedule.name,
			schedule.expr
		FROM schedule
		WHERE schedule.id = $1`

	var schedule Schedule
	dest := []any{
		&schedule.ID,
		&schedule.Name,
		&schedule.Expr,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := r.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Read(id)
		}

		return Schedule{}, err
	}

	return schedule, nil
}

func (r *PostgresRepository) List() ([]Schedule, error) {
	stmt := `
		SELECT
			schedule.id,
			schedule.name,
			schedule.expr
		FROM schedule`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []Schedule
	for rows.Next() {
		var schedule Schedule
		dest := []any{
			&schedule.ID,
			&schedule.Name,
			&schedule.Expr,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, database.ErrRetry) {
				return r.List()
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

func (r *PostgresRepository) Update(schedule Schedule) error {
	return errors.New("TODO: not implemented")
}

func (r *PostgresRepository) Delete(schedule Schedule) error {
	stmt := `
		DELETE FROM schedule
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	err := database.Exec(r.conn, ctx, stmt, schedule.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Delete(schedule)
		}

		return err
	}

	return nil
}
