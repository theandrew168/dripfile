package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/postgres"
)

type transferQueue struct {
	conn    *pgxpool.Pool
	storage database.Storage
}

const (
	StatusNew     = "new"
	StatusRunning = "running"
	StatusSuccess = "success"
	StatusError   = "error"
)

var (
	queryTimeout = 3 * time.Second
)

func NewTransferQueue(conn *pgxpool.Pool, storage database.Storage) *transferQueue {
	q := transferQueue{
		conn:    conn,
		storage: storage,
	}
	return &q
}

// insert transfer ID into the queue table
// https://webapp.io/blog/postgres-is-the-answer/
func (q *transferQueue) Publish(transfer core.Transfer) error {
	stmt := `
		INSERT INTO transfer_queue
			(transfer_id, status)
		VALUES
			($1, $2)`

	args := []interface{}{
		transfer.ID,
		StatusNew,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(q.conn, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return q.Publish(transfer)
		}

		return err
	}

	return nil
}

// atomically claim a job
// https://webapp.io/blog/postgres-is-the-answer/
func (q *transferQueue) Subscribe() (core.Transfer, error) {
	stmt := `
		UPDATE transfer_queue
		SET status = 'running'
		WHERE id = (
			SELECT id
			FROM transfer_queue
			WHERE status = 'new'
			ORDER BY id
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		) RETURNING transfer_id`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	// read the transfer ID that was selected
	var transferID string
	row := q.conn.QueryRow(ctx, stmt)
	err := postgres.Scan(row, &transferID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return q.Subscribe()
		}

		return core.Transfer{}, err
	}

	// lookup the transfer details
	return q.storage.Transfer.Read(transferID)
}
