package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

var (
	queryTimeout = 3 * time.Second
)

func exec(conn *pgxpool.Pool, ctx context.Context, stmt string, args ...interface{}) error {
	_, err := conn.Exec(ctx, stmt, args...)
	if err != nil {
		// check for more specific errors
		// https://github.com/jackc/pgx/wiki/Error-Handling
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// check for stale connections (database restarted)
			if pgErr.Code == pgerrcode.AdminShutdown {
				return core.ErrRetry
			}
		}

		// otherwise bubble the error as-is
		return err
	}

	return nil
}

func scan(row pgx.Row, dest ...interface{}) error {
	err := row.Scan(dest...)
	if err != nil {
		// check for empty result (from QueryRow)
		if errors.Is(err, pgx.ErrNoRows) {
			return core.ErrNotExist
		}

		// check for more specific errors
		// https://github.com/jackc/pgx/wiki/Error-Handling
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// check for duplicate primary keys
			if pgErr.Code == pgerrcode.UniqueViolation {
				return core.ErrExist
			}
			// check for stale connections (database restarted)
			if pgErr.Code == pgerrcode.AdminShutdown {
				return core.ErrRetry
			}
		}

		// otherwise bubble the error as-is
		return err
	}

	return nil
}
