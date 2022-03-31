package migrate

import (
	"context"
	"embed"
	"io/fs"
	"sort"

	"github.com/jackc/pgx/v4"

	"github.com/theandrew168/dripfile/pkg/log"
)

//go:embed migration
var migrationsFS embed.FS

func Migrate(ctx context.Context, conn *pgx.Conn, logger log.Logger) error {
	// create migrations table if it doesn't exist
	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS migration (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)`)
	if err != nil {
		return err
	}

	// get migrations that are already applied
	rows, err := conn.Query(ctx, "SELECT name FROM migration")
	if err != nil {
		return err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return err
		}
		applied[name] = true
	}

	// get migrations that should be applied (from migrations FS)
	subdir, _ := fs.Sub(migrationsFS, "migration")
	migrations, err := fs.ReadDir(subdir, ".")
	if err != nil {
		return err
	}

	// determine missing migrations
	var missing []string
	for _, migration := range migrations {
		name := migration.Name()
		if _, ok := applied[name]; !ok {
			missing = append(missing, name)
		}
	}

	// sort missing migrations to preserve order
	sort.Strings(missing)
	for _, name := range missing {
		logger.Info("applying: %s\n", name)

		// apply the missing ones
		sql, err := fs.ReadFile(subdir, name)
		if err != nil {
			return err
		}
		_, err = conn.Exec(ctx, string(sql))
		if err != nil {
			return err
		}

		// update migrations table
		_, err = conn.Exec(ctx, "INSERT INTO migration (name) VALUES ($1)", name)
		if err != nil {
			return err
		}
	}

	logger.Info("migrations up to date")
	return nil
}
