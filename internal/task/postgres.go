package task

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/postgres"
)

var (
	queryTimeout = 3 * time.Second
)

type postgresQueue struct {
	conn *pgxpool.Pool
}

func NewPostgresQueue(conn *pgxpool.Pool) Queue {
	q := postgresQueue{
		conn: conn,
	}
	return &q
}

// insert transfer ID into the queue table
// https://webapp.io/blog/postgres-is-the-answer/
func (q *postgresQueue) Publish(task Task) error {
	stmt := `
		INSERT INTO task_queue
			(kind, info, status)
		VALUES
			($1, $2, $3)`

	args := []interface{}{
		task.Kind,
		task.Info,
		task.Status,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(q.conn, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return q.Publish(task)
		}

		return err
	}

	return nil
}

// atomically claim a job
// https://webapp.io/blog/postgres-is-the-answer/
func (q *postgresQueue) Subscribe() (Task, error) {
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

	row := q.conn.QueryRow(ctx, stmt)
	err := postgres.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return q.Subscribe()
		}

		return Task{}, err
	}

	return task, nil
}
