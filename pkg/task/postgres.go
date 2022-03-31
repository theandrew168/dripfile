package task

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/postgres"
)

var (
	queryTimeout = 3 * time.Second
)

type postgresQueue struct {
	pool *pgxpool.Pool
}

func NewPostgresQueue(pool *pgxpool.Pool) Queue {
	q := postgresQueue{
		pool: pool,
	}
	return &q
}

// insert transfer ID into the queue table
// https://webapp.io/blog/postgres-is-the-answer/
func (q *postgresQueue) Push(task Task) error {
	stmt := `
		INSERT INTO task_queue
			(kind, info, status, error)
		VALUES
			($1, $2, $3, $4)`

	args := []interface{}{
		task.Kind,
		task.Info,
		task.Status,
		task.Error,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(q.pool, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return q.Push(task)
		}

		return err
	}

	return nil
}

// atomically claim tasks (only run one per worker)
// https://webapp.io/blog/postgres-is-the-answer/
func (q *postgresQueue) Pop() (Task, error) {
	stmt := `
		UPDATE task_queue
		SET status = 'running'
		WHERE id = (
			SELECT id
			FROM task_queue
			WHERE status = 'new'
			ORDER BY id
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		) RETURNING id, kind, info, status`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	// read the task that was selected
	var task Task
	dest := []interface{}{
		&task.ID,
		&task.Kind,
		&task.Info,
		&task.Status,
	}

	row := q.pool.QueryRow(ctx, stmt)
	err := postgres.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return q.Pop()
		}

		return Task{}, err
	}

	return task, nil
}

func (q *postgresQueue) Update(task Task) error {
	stmt := `
		UPDATE task_queue
		SET
			kind = $2,
			info = $3,
			status = $4,
			error = $5
		WHERE id = $1`

	args := []interface{}{
		task.ID,
		task.Kind,
		task.Info,
		task.Status,
		task.Error,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(q.pool, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return q.Update(task)
		}

		return err
	}

	return nil
}