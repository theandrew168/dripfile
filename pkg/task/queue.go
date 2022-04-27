package task

import (
	"context"
	"errors"
	"time"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/database"
)

// default query timeout
var queryTimeout = 3 * time.Second

type Queue struct {
	db database.Conn
}

func NewQueue(db database.Conn) *Queue {
	q := Queue{
		db: db,
	}
	return &q
}

// insert transfer ID into the queue table
// https://webapp.io/blog/database-is-the-answer/
func (q *Queue) Push(task Task) error {
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

	err := database.Exec(q.db, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return q.Push(task)
		}

		return err
	}

	return nil
}

// atomically claim tasks (only run one per worker)
// https://webapp.io/blog/database-is-the-answer/
func (q *Queue) Pop() (Task, error) {
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

	row := q.db.QueryRow(ctx, stmt)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return q.Pop()
		}

		return Task{}, err
	}

	return task, nil
}

func (q *Queue) Update(task Task) error {
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

	err := database.Exec(q.db, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return q.Update(task)
		}

		return err
	}

	return nil
}
