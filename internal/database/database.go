package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const Timeout = 3 * time.Second

var (
	// based the os package error names:
	// https://pkg.go.dev/os#pkg-variables
	ErrExist    = errors.New("database: already exists")
	ErrNotExist = errors.New("database: does not exist")

	// storage errors
	ErrRetry    = errors.New("database: retry storage operation")
	ErrConflict = errors.New("database: conflict in storage operation")
)

// Common interface for pgx.Conn, pgx.Pool, pgx.Tx, etc
// https://github.com/jackc/pgx/issues/875
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

	pool, err := pgxpool.ConnectConfig(ctx, config)
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
			// check for duplicate primary keys
			if pgErr.Code == pgerrcode.UniqueViolation {
				return ErrExist
			}
			// check for stale connections (database restarted)
			if pgErr.Code == pgerrcode.AdminShutdown {
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
			// check for duplicate primary keys
			if pgErr.Code == pgerrcode.UniqueViolation {
				return ErrExist
			}
			// check for stale connections (database restarted)
			if pgErr.Code == pgerrcode.AdminShutdown {
				return ErrRetry
			}
		}

		// otherwise bubble the error as-is
		return err
	}

	return nil
}
