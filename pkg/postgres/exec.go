package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/theandrew168/dripfile/pkg/core"
)

func Exec(pg Interface, ctx context.Context, stmt string, args ...interface{}) error {
	_, err := pg.Exec(ctx, stmt, args...)
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
