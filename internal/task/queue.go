package task

import (
	"context"
	"errors"
	"time"

	"github.com/theandrew168/dripfile/internal/database"
)

// default query timeout
const timeout = 3 * time.Second

type Queue struct {
	db database.Conn
}

func NewQueue(db database.Conn) *Queue {
	q := Queue{
		db: db,
	}
	return &q
}

func (q *Queue) Submit(t Task) error {
	stmt := `
		INSERT INTO task_queue
			(kind, info, status, error)
		VALUES
			($1, $2, $3, $4)`

	args := []any{
		t.Kind,
		t.Info,
		t.Status,
		t.Error,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := database.Exec(q.db, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return q.Submit(t)
		}

		return err
	}

	return nil
}

// atomically claim tasks (only run one per worker)
// https://webapp.io/blog/database-is-the-answer/
func (q *Queue) Claim() (Task, error) {
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
		)
		RETURNING id, kind, info, status`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// read the task that was selected
	var t Task
	dest := []any{
		&t.ID,
		&t.Kind,
		&t.Info,
		&t.Status,
	}

	row := q.db.QueryRow(ctx, stmt)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return q.Claim()
		}

		return Task{}, err
	}

	return t, nil
}

func (q *Queue) Finish(t Task) error {
	if t.Status == StatusFailure {
		return q.finishFailure(t)
	} else {
		return q.finishSuccess(t)
	}
}

func (q *Queue) finishFailure(t Task) error {
	stmt := `
		UPDATE task_queue
		SET
			status = $2,
			error = $3
		WHERE id = $1`

	args := []any{
		t.ID,
		t.Status,
		t.Error,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := database.Exec(q.db, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return q.finishFailure(t)
		}

		return err
	}

	return nil
}

func (q *Queue) finishSuccess(t Task) error {
	stmt := `
		DELETE FROM task_queue
		WHERE id = $1`

	args := []any{
		t.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := database.Exec(q.db, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return q.finishSuccess(t)
		}

		return err
	}

	return nil
}
