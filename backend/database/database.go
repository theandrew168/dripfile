package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*

Possible Errors:

create (exec) - constraint violation (conflict w/ cols)
	UUIDs will always be unique, but other fields could cause issues.
	For example, if adding a user with an email (NOT NULL UNIQUE) that
	already exists, this would cause a constraint violation.

list (collectRows) - none
	Read only. No potential errors here: either returns some rows or none.

read (collectOneRow) - does not exist
	Read only. Only potential error is not finding a row with the specified ID.

update (scan) - does not exist (coalesce to conflict), constraint violation (conflict w/ cols)
	Scan is used here to read the deleted record's version and check for ErrNoRows.
	This is the most complex operation: multiple things could go wrong.
	  1. The record being updated doesn't exist (indistinguishable from TOCTOU check, will appear as conflict)
	  	 This is technically a programming error on the caller's side: updating a record
		 that doesn't exist yet.
	  2. The record being updated causes a constraint violation (dupe values in a UNIQUE column)
	  	 Probably need to communicate this back to the user in one way or another.
	  3. The record being updated was changed between fetch and update (TOCTOU race condition)
	     Based on Alex Edwards' approach to optimistic concurrency control in Let's Go Further.
		 The record exists, but was updated by someone (or something) else before the current
		 request completed. Probably need to tell the user to try again.

delete (scan) - does not exist
	Scan is used here to read the deleted record's ID and check for ErrNoRows.
	Only potential error is not finding a row with the specified ID. This could just
	ignore cases where the ID doesn't exist (and nothing gets deleted) but I think it is
	better UX / DX to _know_ if the delete was successful (204) vs no record was deleted (404).

*/

const Timeout = 3 * time.Second

var (
	// based the os package error names:
	// https://pkg.go.dev/os#pkg-variables
	ErrExist    = errors.New("database: already exists")
	ErrNotExist = errors.New("database: does not exist")

	// storage errors
	ErrRetry    = errors.New("database: retry storage operation")
	ErrConflict = errors.New("database: conflict in storage operation")

	// data errors
	ErrInvalidUUID = errors.New("database: invalid UUID")
)

// Common interface for pgx.Conn, pgx.Pool, pgx.Tx, etc
// https://github.com/jackc/pgx/issues/644
type Conn interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func Connect(databaseURI string) (*pgx.Conn, error) {
	ctx := context.Background()

	config, err := pgx.ParseConfig(databaseURI)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		conn.Close(ctx)
		return nil, err
	}

	return conn, nil
}

func ConnectPool(databaseURI string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(databaseURI)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func Exec(db Conn, ctx context.Context, stmt string, args ...any) error {
	_, err := db.Exec(ctx, stmt, args...)
	if err != nil {
		// check for more specific errors
		// https://github.com/jackc/pgx/wiki/Error-Handling
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// check for constraint violations
			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return ErrConflict
			}
			// check for stale connections (database restarted)
			if pgerrcode.IsOperatorIntervention(pgErr.Code) {
				return ErrRetry
			}
		}

		// otherwise bubble the error as-is
		return err
	}

	return nil
}

func Scan(row pgx.Row, dest ...any) error {
	err := row.Scan(dest...)
	if err != nil {
		// check for empty result (from QueryRow)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotExist
		}

		// check for more specific errors
		// https://github.com/jackc/pgx/wiki/Error-Handling
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// check for other constraint violations
			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return ErrConflict
			}
			// check for stale connections (database restarted)
			if pgerrcode.IsOperatorIntervention(pgErr.Code) {
				return ErrRetry
			}
		}

		// otherwise bubble the error as-is
		return err
	}

	return nil
}
