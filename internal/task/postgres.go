package task

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/postgres"
)

var (
	queryTimeout = 3 * time.Second
)

type postgresQueue struct {
	conn *pgx.Conn
	pool *pgxpool.Pool
}

func NewPostgresQueue(conn *pgx.Conn, pool *pgxpool.Pool) Queue {
	q := postgresQueue{
		conn: conn,
		pool: pool,
	}
	return &q
}

// insert transfer ID into the queue table
// https://webapp.io/blog/postgres-is-the-answer/
func (q *postgresQueue) Push(task Task) error {
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

func (q *postgresQueue) Listen(ctx context.Context) error {
	// listen on the task queue status channel
	_, err := q.conn.Exec(ctx, "listen task_queue_status_channel")
	if err != nil {
		return err
	}

	// block until trigger or ctx deadlines expires
	_, err = q.conn.WaitForNotification(ctx)
	if err != nil {
		return err
	}

	// ready to pop!
	return nil
}
